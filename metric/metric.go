package metric

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Metric struct {
	Id         int
	Device_id  int
	Metric_1   int
	Metric_2   int
	Metric_3   int
	Metric_4   int
	Metric_5   int
	Local_time time.Time
}

func GenerateMetric(deviceId int) error {
	for {
		//imitating ping
		timestep := time.Duration(250 + rand.Intn(5000))
		time.Sleep(timestep * time.Millisecond)
		metricBody, err := json.Marshal(Metric{
			Id:         rand.Intn(10000),
			Device_id:  deviceId,
			Metric_1:   rand.Intn(1000),
			Metric_2:   rand.Intn(1000),
			Metric_3:   rand.Intn(1000),
			Metric_4:   rand.Intn(1000),
			Metric_5:   rand.Intn(1000),
			Local_time: time.Now(),
		})
		if err != nil {
			log.Println("GenerateMetric json.Marshal error: ", err)
			return err
		}

		_, err = http.Post("http://localhost:3000/metric", "application/json", bytes.NewReader(metricBody))
		if err != nil {
			log.Println("GenerateMetric json.Marshal error: ", err)
			return err
		}

		//return nil

	}

}
