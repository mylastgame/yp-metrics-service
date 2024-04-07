package storage

type Storage interface {
	SaveGauge(string, float64)
	SaveCounter(string, int64)
	ResetCounters()
	GetGauges() map[string]float64
	GetCounters() map[string]int64
}
