package redisAlert

import (
	"log"

	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
)

var rpool *pool.Pool

// SetRedisPool Establish a pool of 500 connections to the Redis server listening on
// port repisPort of the local machine.
func SetRedisPool(repisPort string) (err error) {

	rpool, err = pool.New("tcp", "localhost:"+repisPort, 500)
	if err != nil {
		log.Println("GetRedisPool pool.New error:", err)
		return err
	}

	return nil
}

// SaveAlertRedis saves the last notification about bad metric in Redis in map DeviceID_ - Message
func SaveAlertRedis(id int, msg string) error {
	log.Println("SAR DEBUG1")
	conn, err := rpool.Get()
	if err != nil {
		log.Println("SaveAlertRedis  db.Get() error:", err)
		return err
	}

	defer rpool.Put(conn)

	//HSET addhash DeviseID_id, Message:msg
	err = conn.Cmd("HMSET", "DeviceID_"+strconv.Itoa(id), "Message", msg).Err
	if err != nil {
		log.Println("SaveAlertRedis conn.Cmd error:", err)
		return err
	}
	return nil
}
