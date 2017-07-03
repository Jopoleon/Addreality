package main

import (
	"log"
	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
)

var rpool *pool.Pool
var err1 error

func main() {

	//var rpool *pool.Pool
	repisPort := "6379"
	//var err error
	// Establish a pool of 10 connections to the Redis server listening on
	// port 6379 of the local machine.
	rpool, err1 = pool.New("tcp", "localhost:"+repisPort, 10)
	if err1 != nil {
		log.Panic(" pool.New error:", err1)
		return
	}
	//conn, err := rpool.Get()
	//if err != nil {
	//	log.Fatalln("SaveAlertRedis  db.Get() error:", err)
	//	return
	//}
	//
	//defer rpool.Put(conn)
	var id = 21231
	var msg = "2131ssdaASDXCZX "
	//HSET addhash DeviseID:id, Message:msg
	err := SaveAlertRedis(id, msg)
	if err != nil {
		log.Fatal("Save RedisA lert err: ", err)
	}
	//err = conn.Cmd("HMSET", "DeviseID:"+strconv.Itoa(id), "Message", msg).Err
	//if err != nil {
	//	log.Fatalln("SaveAlertRedis conn.Cmd error:", err)
	//	return
	//}
	CheckResult(id)
}

func SaveAlertRedis(id int, msg string) error {
	log.Println("SAR DEBUG1")
	conn, err := rpool.Get()
	if err != nil {
		log.Fatalln("SaveAlertRedis  db.Get() error:", err)
		return err
	}
	log.Println("SAR DEBUG2")
	defer rpool.Put(conn)

	//HSET addhash DeviseID:id, Message:msg
	err = conn.Cmd("HMSET", "DeviceID_"+strconv.Itoa(id), "Message", msg).Err
	if err != nil {
		log.Fatalln("SaveAlertRedis conn.Cmd error:", err)
		return err
	}
	log.Println("SAR DEBUG3")
	return nil
}

func CheckResult(id int) {
	conn, err := rpool.Get()
	if err != nil {
		log.Fatalln("SaveAlertRedis  db.Get() error:", err)
		return
	}
	price, err := conn.Cmd("HGET", "DeviceID_"+strconv.Itoa(id), "Message").Str()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Result form redis:", price)
	return
}
