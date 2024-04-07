package counter

type Repo interface {
	Save(*Counter) error
	Get(string) (*Counter, bool)
}
