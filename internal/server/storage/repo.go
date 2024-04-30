package storage

type Repo interface {
	Set(string, string, string) error
	Get(string, string) (string, error)
	GetCounters() []string
	GetGauges() []string
	SetGauge(string, float64)
	SetCounter(string, int64)
	GetGauge(string) (float64, bool)
	GetCounter(string) (int64, bool)
}
