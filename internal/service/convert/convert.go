package convert

import "strconv"

func GaugeToString(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func CounterToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

func StringToGauge(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func StringToCounter(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
