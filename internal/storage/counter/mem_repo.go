package counter

import (
	"fmt"
	counterDomain "github.com/mylastgame/yp-metrics-service/internal/domain/counter"
)

type MemRepo struct {
	storage map[string]counterDomain.ValType
}

func NewMemRepo() *MemRepo {
	return &MemRepo{storage: make(map[string]counterDomain.ValType)}
}

func (r *MemRepo) Save(c *counterDomain.Counter) error {
	_, ok := r.storage[c.Title]
	if !ok {
		r.storage[c.Title] = c.Val
	} else {
		r.storage[c.Title] += c.Val
	}

	fmt.Println("counter storage: ", r.storage)
	return nil
}

func (r *MemRepo) Get(title string) (*counterDomain.Counter, bool) {
	v, ok := r.storage[title]

	if ok {
		return counterDomain.New(title, v), ok
	} else {
		return &counterDomain.Counter{}, ok
	}
}
