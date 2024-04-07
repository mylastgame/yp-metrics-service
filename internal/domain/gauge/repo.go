package gauge

type Repo interface {
	Save(*Gauge) error
	Get(string) (*Gauge, bool)
}
