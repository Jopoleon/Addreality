package metric

import (
	"errors"
	"strconv"
	"time"

	"github.com/Jopoleon/AddRealtyTask/config"
)

type Metric struct {
	Device_id  int
	Metric_1   int
	Metric_2   int
	Metric_3   int
	Metric_4   int
	Metric_5   int
	Local_time time.Time
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
