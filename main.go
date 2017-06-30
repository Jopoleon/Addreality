package main

import (
	"log"
	"net/http"

	"flag"

	"database/sql"

	"encoding/json"
	"io/ioutil"

	"strconv"

	"errors"

	"fmt"

	"github.com/Jopoleon/AddRealtyTask/config"
	"github.com/Jopoleon/AddRealtyTask/db"
	"github.com/Jopoleon/AddRealtyTask/metric"
	"github.com/gorilla/context"
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

	log.Printf("CONFIG FILE MAIN: %+v", Config)
}
func main() {
	for i := 0; i < 50; i++ {
		go metric.GenerateMetric(i)
	}
	DB, dberr = db.SetDB(Config.Host, Config.Port, Config.User, Config.Password, Config.DBname)
	if dberr != nil {
		log.Fatalln("main() db.SetDB err: ", dberr)
		return
	}

	http.HandleFunc("/metric", MetricHandler)
	err := http.ListenAndServe(":"+Config.ServerPort, context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatalln(err)
		return
	}
}

//CheckMetricValues checks if given metric meet the conditions given in config file
func CheckMetricValues(m metric.Metric, conf *config.ConfigType) (ok bool, met metric.Metric, err error) {

	if conf.Metric_1_Max < m.Metric_1 || m.Metric_1 < conf.Metric_1_Min {
		return false, m, errors.New("Metric 1 of device" + strconv.Itoa(m.Device_id) + " is out of boundaries.")
	}
	if conf.Metric_2_Max < m.Metric_2 || m.Metric_2 < conf.Metric_2_Min {
		return false, m, errors.New("Metric 2 of device" + strconv.Itoa(m.Device_id) + " is out of boundaries.")
	}
	if conf.Metric_3_Max < m.Metric_3 || m.Metric_3 < conf.Metric_3_Min {
		return false, m, errors.New("Metric 3 of device" + strconv.Itoa(m.Device_id) + " is out of boundaries.")
	}
	if conf.Metric_4_Max < m.Metric_4 || m.Metric_4 < conf.Metric_4_Min {
		return false, m, errors.New("Metric 4 of device" + strconv.Itoa(m.Device_id) + " is out of boundaries.")
	}
	if conf.Metric_5_Max < m.Metric_5 || m.Metric_5 < conf.Metric_5_Min {
		return false, m, errors.New("Metric 5 of device" + strconv.Itoa(m.Device_id) + " is out of boundaries.")
	}
	return true, m, nil
}

func MetricHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("MetricHandler ioutil.ReadAll(resp.Body) error", err)
		return
	}
	var metricData metric.Metric
	err = json.Unmarshal(body, &metricData)
	if err != nil {
		log.Println("MetricHandler json.Unmarshal error", err)
		return
	}
	ok, met, err := CheckMetricValues(metricData, Config)
	if !ok {
		log.Printf("Bad metric: \n %s, %+v", err, met)
		//save alert to PostgreSQL
		//fmt.Sprintf("%+v",met)
		err = db.SaveAlert(met, err.Error()+fmt.Sprintf("%+v", met), DB)
		if err != nil {
			log.Println("MetricHandler db.SaveAlert error", err)
			return
		}
		//send alert email here
		//save alert to Redis
		return
	}
	err = db.SaveMetric(metricData, DB)
	if err != nil {
		log.Println("MetricHandler db.SaveMetric error", err)
		return
	}
}
