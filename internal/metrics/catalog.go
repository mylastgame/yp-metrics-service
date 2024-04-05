package metrics

import (
	storage "github.com/mylastgame/yp-metrics-service/internal/storage"
	"strconv"
)

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

//type Catalog struct {
//	gauge []string
//	counter []string
//}
//
//var Metrics = Catalog{
//	gauge: []string{"gauge1", "gauge2", "gauge3", "gauge4", "gauge5"}
//	counter: []string{"counter1", "counter2", "counter3", "counter4", "counter5"}
//}

var metrics = map[string][]string{
	TypeGauge:   []string{"gauge1", "gauge2", "gauge3", "gauge4", "gauge5"},
	TypeCounter: []string{"counter1", "counter2", "counter3", "counter4", "counter5"},
}

func TypeExists(mtype string) bool {
	_, ok := metrics[mtype]
	if !ok {
		return false
	}

	return true
}

func MetricExists(mtype, title string) bool {
	titles, ok := metrics[mtype]
	if !ok {
		return false
	}

	for _, mtitle := range titles {
		if mtitle == title {
			return true
		}
	}

	return false
}

func Save(mtype, key, v string) (bool, error) {
	if mtype == TypeGauge {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false, err
		}
		storage.Storage.AddGauge(key, value)
	}

	if mtype == TypeCounter {
		value, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return false, err
		}
		storage.Storage.AddCounter(key, value)
	}

	return true, nil
}
