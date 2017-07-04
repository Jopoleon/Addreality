package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

type User struct {
	ID    int
	Name  string
	Email string
}
type Device struct {
	ID     int
	Name   string
	UserID int
}

func initPostgresTables(DB *sql.DB) error {

	err := initUsersTable(DB)
	if err != nil {
		log.Println("initPostgreTables initUsersTable(DB) error: ", err)
		return err
	}
	err = initDeviceTable(DB)
	if err != nil {
		log.Println("initPostgreTables initDeviceTable(DB) error: ", err)
		return err
	}
	err = initDeviceMetricsTable(DB)
	if err != nil {
		log.Println("initPostgreTables initDeviceMetricsTable(DB) error: ", err)
		return err
	}
	err = initDeviceAlertsTable(DB)
	if err != nil {
		log.Println("initPostgreTables initDeviceAlertsTable(DB) error: ", err)
		return err
	}
	return nil

}
func initUsersTable(DB *sql.DB) error {
	res, err := DB.Exec(`CREATE TABLE IF NOT EXISTS ` + ` users
(
  id SERIAL PRIMARY KEY,
  name varchar(255),
  email varchar(255) UNIQUE NOT NULL
)`)
	if err != nil {
		log.Println("initPostgreTables users db.Exec error: ", err)
		return err
	}
	log.Println("Creation of users table is ok: ", res)
	return nil
}
func initDeviceTable(DB *sql.DB) error {
	res, err := DB.Exec(`CREATE TABLE IF NOT EXISTS ` + ` devices
(
  id INT PRIMARY KEY,
  name varchar(255) NOT NULL,
  user_id INT NOT NULL,

  CONSTRAINT devices_user_id_fk FOREIGN KEY(user_id) REFERENCES users (id) ON DELETE CASCADE
)`)
	if err != nil {
		log.Println("initPostgreTables devices db.Exec error: ", err)
		return err
	}
	log.Println("Creation of devices table is ok: ", res)
	return nil
}
func initDeviceMetricsTable(DB *sql.DB) error {
	res, err := DB.Exec(`CREATE TABLE IF NOT EXISTS ` + `device_metrics
(
    id SERIAL PRIMARY KEY,
    device_id INT NOT NULL,
    metric_1 INT,
    metric_2 INT,
    metric_3 INT,
    metric_4 INT,
    metric_5 INT,
    local_time TIMESTAMP,
    server_time TIMESTAMP DEFAULT NOW(),

    CONSTRAINT device_metrics_device_id_fk FOREIGN KEY (device_id) REFERENCES devices (id) ON DELETE CASCADE
)`)
	if err != nil {
		log.Println("initPostgreTables device_metrics db.Exec error: ", err)
		return err
	}
	log.Println("Creation of device_metrics table is ok: ", res)
	return nil
}
func initDeviceAlertsTable(DB *sql.DB) error {
	res, err := DB.Exec(`CREATE TABLE IF NOT EXISTS ` + `device_alerts
(
  id SERIAL PRIMARY KEY,
  device_id INT,
  message TEXT
)`)
	if err != nil {
		log.Println("initPostgreTables device_alerts db.Exec error: ", err)
		return err
	}
	log.Println("Creation of device_alerts table is ok: ", res)

	return nil
}

func addUser(DB *sql.DB) (userID int, err error) {

	sqlStatement := `
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING id`
	id := 0
	err = DB.QueryRow(sqlStatement, "TestUser", "egortictac3@gmail.com").Scan(&id)
	if err != nil {
		log.Fatalln("initUsers() db.QueryRow error: ", err)
		return id, err
	}
	fmt.Println("New User ID is:", id)
	return id, nil
}
func initDevices(DB *sql.DB, userID int) error {
	sqlStatement := `
INSERT INTO devices (id, name, user_id)
VALUES ($1, $2, $3)
RETURNING id`
	for i := 0; i < 10000; i++ {
		id := 0
		err := DB.QueryRow(sqlStatement, i, "Device"+strconv.Itoa(i), userID).Scan(&id)
		if err != nil {
			log.Fatalln("initDevices db.QueryRow error: ", err)
			return err
		}
	}
	return nil
}
