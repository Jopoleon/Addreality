package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"strconv"

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
	var uID int
	log.Println("Length of DEVICES table is: ", countDevices)
	if countUsers == 0 {
		userID, err := addUser(DB)
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
		uID = userID
		if countDevices == 0 {
			err = initDevices(DB, userID)
			if err != nil {
				log.Fatalln(err)
				return nil, err
			}
		}
	}
	if countDevices == 0 {
		err = initDevices(DB, uID)
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
	}
	return DB, nil
}

func SaveAlert(m metric.Metric, msg string, DB *sql.DB) error {

	sqlStatement := `
INSERT INTO device_alerts (device_id,message)
VALUES ($1, $2)
RETURNING id`
	var id int
	err := DB.QueryRow(sqlStatement, m.Device_id, msg).Scan(&id)
	if err != nil {
		log.Fatalln("SaveAlert db.QueryRow error: ", err)
		return err
	}
	return nil
}

func SaveMetric(m metric.Metric, DB *sql.DB) error {
	sqlStatement := `
INSERT INTO device_metrics (device_id,metric_1,metric_2,metric_3,metric_4,metric_5,local_time,server_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id`
	var d int
	err := DB.QueryRow(sqlStatement, m.Device_id, m.Metric_1, m.Metric_2, m.Metric_3, m.Metric_4, m.Metric_5, m.Local_time, time.Now()).Scan(&d)
	if err != nil {
		log.Fatalln("SaveMetric device_metrics db.QueryRow error: ", err)
		return err
	}
	return nil
}
func GetUserInfo(deviceID int, DB *sql.DB) (user User, err error) {
	sqlStatement2 := `SELECT * FROM devices WHERE name=$1;`
	var device Device
	row := DB.QueryRow(sqlStatement2, "Device"+strconv.Itoa(deviceID))
	err = row.Scan(&device.ID, &device.Name, &device.UserID)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, sql.ErrNoRows
	case nil:
		fmt.Println(device)
		sqlStatement2 := `SELECT * FROM users WHERE id=$1;`
		//var user1 User
		row := DB.QueryRow(sqlStatement2, device.UserID)
		err := row.Scan(&user.ID, &user.Name, &user.Email)
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return user, err
		case nil:
			fmt.Println("GetUserInfo user found: ", user)
			return user, err
		default:
			panic("GetUserInfo SELECT * FROM users row.Scan error: " + err.Error())
		}
	default:
		panic("GetUserInfo SELECT * FROM devices row.Scan error: " + err.Error())
	}
	//return nil, nil
}
