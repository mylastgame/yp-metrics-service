package counter

type Repo interface {
	Add(*Counter) error
	Get(string) (*Counter, bool)
	GetAll() []*Counter
}
