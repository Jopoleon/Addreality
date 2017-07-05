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
	"github.com/Jopoleon/AddRealtyTask/redisAlert"
	"github.com/Jopoleon/AddRealtyTask/sendemail"
)

var (
	Config    *config.ConfigType
	DB        *sql.DB
	dberr     error
	EmailAuth sendemail.SmtpAuth
	Sendmail  *bool
)

func init() {
	var configFileName string

	flag.StringVar(&configFileName, "config", "config.json",
		"Specify configuration file name to use. File should be in folder you starting the application")
	Sendmail = flag.Bool("sendmail", false,
		"Specify sending emails with notifications about bad metrics or not")

	//flag.StringVar(&Sendmail, )
	flag.Parse()

	config.InitConf(configFileName)

	Config = config.GetConfig()
	EmailAuth = sendemail.AuthMailBox(sendemail.EmailUser{Config.EmailLogin, Config.EmailPassword, Config.EmailServer, Config.EmailPort})

	fmt.Printf("%s %+v \n", "CONFIG FILE MAIN:", Config)
	log.Printf("%s %s", "Sending mail with notifications: ", *Sendmail)
}

func main() {
	// starting gorutines which imitating metrics from devices
	for i := 0; i < 1; i++ {
		go metric.GenerateMetric(i, Config.ServerPort)
	}
	//setting connetcion to PostgreSQl
	DB, dberr = db.SetDB(Config)
	if dberr != nil {
		log.Println("main() db.SetDB err: ", dberr)
		return
	}
	//setting Redis connection pool
	err := redisAlert.SetRedisPool(Config.RedisPort)
	if dberr != nil {
		log.Fatalln("main() redis.SetRedisPool err: ", dberr)
		return
	}

	defer DB.Close()

	http.HandleFunc("/metric", MetricHandler)
	err = http.ListenAndServe(":"+Config.ServerPort, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Server started on port ", Config.ServerPort)

}

//MetricHandler handels requests from GenerateMetric function on /metric endpoint.
func MetricHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("MetricHandler ioutil.ReadAll(resp.Body) error", err)
		return
	}

	var metricData []metric.Metric
	err = json.Unmarshal(body, &metricData)
	if err != nil {
		log.Println("MetricHandler json.Unmarshal error", err)
		return
	}
	defer r.Body.Close()
	// iterating through all metrics
	for _, md := range metricData {
		ok, met, err1 := metric.CheckMetricValues(md, Config)
		if !ok {
			msg := err1.Error() + fmt.Sprintf("%+v", met)

			log.Printf("Bad metric: \n %s, %+v", msg, met)
			//saving alert to PostgreSQL
			err = db.SaveAlert(met, msg, DB)
			if err != nil {
				log.Println("MetricHandler db.SaveAlert error", err)
				return
			}

			//saving alert to Redis
			err = redisAlert.SaveAlertRedis(met.Device_id, msg)
			//err = redisAlert.SaveAlertRedis(1233, "TEST MESSAGE")
			if err != nil {
				log.Println("MetricHandler redis.SaveAlertRedis error", err)
				return
			}

			//getting user's email whose device posted bad metric
			userinfo, err := db.GetUserInfo(met.Device_id, DB)
			if err != nil {
				log.Println("MetricHandler db.GetUserInfo error", err)
				return
			}

			//sending alert email here
			if *Sendmail {
				err = sendemail.SendEmailwithMessage(userinfo.Email, msg, EmailAuth)
				if err != nil {
					log.Println("MetricHandler sendemail.SendEmailwithMessage error", err)
					return
				}
			}

			return
		}
		err = db.SaveMetric(md, DB)

		if err != nil {
			log.Println("MetricHandler saving good metric db.SaveMetric error", err)
			return
		}
	}

}
