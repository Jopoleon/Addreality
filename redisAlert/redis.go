package redisAlert

import (
	"log"

	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
)

var rpool *pool.Pool

func SetRedisPool(repisPort string) (err error) {
	//var err error
	// Establish a pool of 10 connections to the Redis server listening on
	// port 6379 of the local machine.
	rpool, err = pool.New("tcp", "localhost:"+repisPort, 500)
	if err != nil {
		log.Panic("GetRedisPool pool.New error:", err)
		return err
	}

	return nil
}

//(
//id INT PRIMARY KEY,
//device_id INT,
//message TEXT
//)
type Alert struct {
	DeviseID string
	Message  string
}

func SaveAlertRedis(id int, msg string) error {
	log.Println("SAR DEBUG1")
	conn, err := rpool.Get()
	if err != nil {
		log.Fatalln("SaveAlertRedis  db.Get() error:", err)
		return err
	}

	defer rpool.Put(conn)

	//HSET addhash DeviseID:id, Message:msg
	err = conn.Cmd("HMSET", "DeviceID_"+strconv.Itoa(id), "Message", msg).Err
	if err != nil {
		log.Fatalln("SaveAlertRedis conn.Cmd error:", err)
		return err
	}
	return nil
}
