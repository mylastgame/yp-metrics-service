package counter

import (
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/counter"
)

type MemRepo struct {
	storage map[string]counter.ValType
}

func NewMemRepo() *MemRepo {
	return &MemRepo{storage: make(map[string]counter.ValType)}
}

func (r *MemRepo) Add(c *counter.Counter) error {
	_, ok := r.storage[c.Title]
	if !ok {
		r.storage[c.Title] = c.Val
	} else {
		r.storage[c.Title] += c.Val
	}

	fmt.Println("counter storage: ", r.storage)
	return nil
}

func (r *MemRepo) Get(title string) (*counter.Counter, bool) {
	v, ok := r.storage[title]

	if ok {
		return counter.New(title, v), ok
	} else {
		return &counter.Counter{}, ok
	}
}

func (r *MemRepo) GetAll() []*counter.Counter {
	res := make([]*counter.Counter, 0)

	for t, v := range r.storage {
		res = append(res, counter.New(t, counter.ValType(v)))
	}

	return res
}
