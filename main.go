package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Jopoleon/AddRealtyTask/config"
	"github.com/Jopoleon/AddRealtyTask/db"
	"github.com/Jopoleon/AddRealtyTask/metric"
	"github.com/Jopoleon/AddRealtyTask/redis"
	"github.com/Jopoleon/AddRealtyTask/sendemail"
)

var (
	Config *config.ConfigType
	DB     *sql.DB
	dberr  error
)

func init() {
	var configFileName string

	flag.StringVar(&configFileName, "config", "config.json",
		"Specify configuration file name to use. File should be in folder you starting the application")

	flag.Parse()

	config.InitConf(configFileName)

	Config = config.GetConfig()
	sendemail.AuthMailBox(sendemail.EmailUser{Config.EmailLogin, Config.EmailPassword, Config.EmailServer, Config.EmailPort})

	log.Printf("CONFIG FILE MAIN: %+v", Config)
}
func main() {
	for i := 0; i < 10000; i++ {
		go metric.GenerateMetric(i, Config.ServerPort)
	}

	DB, dberr = db.SetDB(Config.Host, Config.Port, Config.User, Config.Password, Config.DBname)
	if dberr != nil {
		log.Fatalln("main() db.SetDB err: ", dberr)
		return
	}
	_, err := redis.SetRedisPool(Config.RedisPort)
	if dberr != nil {
		log.Fatalln("main() redis.SetRedisPool err: ", dberr)
		return
	}
	//defer DBCloser()
	defer DBCloser()
	log.Println("Debug1")
	http.HandleFunc("/metric", MetricHandler)
	log.Println("Debug2")
	log.Println("Debug3")
	//err := http.ListenAndServe(":"+Config.ServerPort, context.ClearHandler(http.DefaultServeMux))
	err = http.ListenAndServe(":"+Config.ServerPort, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

}
func DBCloser() {
	log.Println("[WARNING] DB connetcion closed.")
	DB.Close()
}

func MetricHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("MetricHandler ioutil.ReadAll(resp.Body) error", err)
		return
	}
	var metricData []metric.Metric
	err1 := json.Unmarshal(body, &metricData)
	if err != nil {
		log.Println("MetricHandler json.Unmarshal error", err1)
		return
	}
	log.Println("Recieved metric data: ", metricData)
	for _, md := range metricData {
		ok, met, err := metric.CheckMetricValues(md, Config)
		if !ok {

			log.Printf("Bad metric: \n %s, %+v", err, met)
			//save alert to PostgreSQL
			err = db.SaveAlert(met, err.Error()+fmt.Sprintf("%+v", met), DB)
			if err != nil {
				log.Println("MetricHandler db.SaveAlert error", err)
				return
			}
			//save alert to Redis
			err = redis.SaveAlertRedis(md.Device_id, err.Error()+fmt.Sprintf("%+v", met))
			if err != nil {
				log.Println("MetricHandler redis.SaveAlertRedis error", err)
				return
			}
			//sav
			//send alert email here

			return
		}
		err = db.SaveMetric(md, DB)

		if err != nil {
			log.Println("MetricHandler db.SaveMetric error", err)
			return
		}
	}

}
