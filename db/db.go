package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Jopoleon/AddRealtyTask/metric"
	_ "github.com/lib/pq"
)

var (
	DB    *sql.DB
	dberr error
)

func SetDB(host, port, user, password, dbname string) (DB *sql.DB, err error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s  sslmode=disable",
		host, port, user, dbname, password)
	DB, dberr = sql.Open("postgres", psqlInfo)

	if dberr != nil {
		log.Fatalln("SetDB sql.Open error: ", dberr)
		return nil, dberr
	}
	err = DB.Ping()
	if err != nil {
		log.Fatalln("SetDB DB.Ping() error: ", err)
		return nil, err
	}
	log.Println("Successfully connected to " + dbname + " database")
	//res, err := DB.Exec(`CREATE DATABASE ` + dbname)
	//if err.Error() != "pq: база данных \"postgres\" уже существует" {
	//	log.Println("SetDB DB.Exec error: CREATE DATABASE IF NOT EXISTS ", dbname, " ", err)
	//	return nil, err
	//}
	//log.Println("Creation of DATABASE is ok: ", res)

	err = initPostgresTables(DB)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	var countUsers int
	err = DB.QueryRow(`SELECT COUNT(*) as count FROM users`).Scan(&countUsers)
	if err != nil {
		log.Fatalln("SetDB DB.Query(`SELECT count(*) error: ", err)
		return nil, err
	}
	log.Println("Length of USERS table is: ", countUsers)
	var countDevices int
	err = DB.QueryRow(`SELECT COUNT(*) as count FROM devices`).Scan(&countDevices)
	if err != nil {
		log.Fatalln("SetDB DB.Query(`SELECT count(*) error: ", err)
		return nil, err
	}
	log.Println("Length of DEVICES table is: ", countDevices)
	if countUsers == 0 {
		userID, err := addUser(DB)
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
		if countDevices == 0 {
			err = initDevices(DB, userID)
			if err != nil {
				log.Fatalln(err)
				return nil, err
			}
		}
	}
	return DB, nil
}
func SaveAlert(m metric.Metric, msg string, DB *sql.DB) error {
	err := DB.Ping()
	if err != nil {
		log.Fatalln("SaveAlert DB.Ping() error: ", err)
		return err
	}

	//device_alerts
	//(
	//	id INT PRIMARY KEY,
	//	device_id INT,
	//	message TEXT
	//)
	sqlStatement := `
INSERT INTO device_alerts (device_id,message)
VALUES ($1, $2)
RETURNING id`
	id := 0
	err = DB.QueryRow(sqlStatement, m.Device_id, msg).Scan(&id)
	if err != nil {
		log.Fatalln("SaveAlert db.QueryRow error: ", err)
		return err
	}
	return nil
}

func SaveMetric(m metric.Metric, DB *sql.DB) error {
	//id INT PRIMARY KEY,
	//	device_id INT NOT NULL,
	//	metric_1 INT,
	//	metric_2 INT,
	//	metric_3 INT,
	//	metric_4 INT,
	//	metric_5 INT,
	//	local_time TIMESTAMP, —Время метрик на устройстве
	//server_time TIMESTAMP DEFAULT NOW() — Серверное время сохранения метрик
	//
	//CONSTRAINT device_metrics_device_id_fk FOREIGN KEY (device_id) REFERENCES devices (id) ON DELETE CASCADE
	err := DB.Ping()
	if err != nil {
		log.Fatalln("SaveMetric DB.Ping() error: ", err)
		return err
	}
	//defer DB.Close()
	sqlStatement := `
INSERT INTO device_metrics (device_id,metric_1,metric_2,metric_3,metric_4,metric_5,local_time,server_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id`
	var d int
	err = DB.QueryRow(sqlStatement, m.Device_id, m.Metric_1, m.Metric_2, m.Metric_3, m.Metric_4, m.Metric_5, m.Local_time, time.Now()).Scan(&d)
	if err != nil {
		log.Fatalln("SaveMetric device_metrics db.QueryRow error: ", err)
		return err
	}
	return nil
}
