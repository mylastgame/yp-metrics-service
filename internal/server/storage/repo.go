package storage

type Repo interface {
	Set(string, string, string) error
	Get(string, string) (string, error)
	GetCounters() []string
	GetGauges() []string
}
