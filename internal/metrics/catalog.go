package metrics

import "strconv"

type Gauge float64
type Counter int64

const (
	GaugeKey   = "gauge"
	CounterKey = "counter"
)

var metricsTypes = [2]string{
	GaugeKey,
	CounterKey,
}

// check if type exists
func TypeExists(mtype string) bool {
	for _, mt := range metricsTypes {
		if mt == mtype {
			return true
		}
	}
	return false
}

// Try to convert string to Gauge type
func ConvertToGauge(value string) (Gauge, error) {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return Gauge(result), nil
}

// Try to convert string to Counter type
func ConvertToCounter(value string) (Counter, error) {
	result, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return 0, err
	}

	return Counter(result), nil
}
