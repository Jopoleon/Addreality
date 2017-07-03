package metric

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"bytes"

	"github.com/Jopoleon/AddRealtyTask/config"
)

type Metric struct {
	//Id         int
	Device_id  int
	Metric_1   int
	Metric_2   int
	Metric_3   int
	Metric_4   int
	Metric_5   int
	Local_time time.Time
}

func GenerateMetric(deviceId int, port string) error {
	for {
		//imitating ping
		timestep := time.Duration(250 + rand.Intn(5000))
		time.Sleep(timestep * time.Millisecond)

		//creating random amount of metrics per request from device
		amount := 1 + rand.Intn(2)
		var mslice []Metric
		for i := 0; i < amount; i++ {
			mslice = append(mslice, Metric{
				//Id:        rand.Intn(10000),
				Device_id: deviceId,
				Metric_1:  119 + rand.Intn(680),
				//	"Metric_1_Max": 680,
				//	"Metric_1_Min": 120,
				Metric_2: 12 + rand.Intn(899),
				//	"Metric_2_Max": 900,
				//	"Metric_2_Min": 12,
				Metric_3: rand.Intn(799),
				//	"Metric_3_Max": 800,
				//	"Metric_3_Min": 46,
				Metric_4: 30 + rand.Intn(699),
				//	"Metric_4_Max": 700,
				//	"Metric_4_Min": 30,
				Metric_5: 200 + rand.Intn(899),
				//	"Metric_5_Max": 900,
				//	"Metric_5_Min": 200
				Local_time: time.Now(),
			})
		}
		metricBody, err := json.Marshal(mslice)
		if err != nil {
			log.Println("GenerateMetric json.Marshal error: ", err)
			return err
		}
		log.Printf("\n thit is how it looks like::::.... %+v", string(metricBody))
		//_, err = http.Post("http://localhost:"+port+"/metric", "application/json", bytes.NewReader(metricBody))
		client := &http.Client{}
		req, err := http.NewRequest("http://localhost:"+port+"/metric", "application/json", bytes.NewReader(metricBody))

		// NOTE this !!
		req.Close = true

		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("user", "pass")
		resp, err := client.Do(req)

		if err != nil {
			log.Println("GenerateMetric json.Marshal error: ", err)
			return err
		}
		if resp != nil {
			defer resp.Body.Close()
		}
	}
}

//CheckMetricValues checks if given metric meets the conditions given in config file
func CheckMetricValues(m Metric, conf *config.ConfigType) (ok bool, met Metric, err error) {

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
