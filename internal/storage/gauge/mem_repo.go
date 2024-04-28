package gauge

import (
	"github.com/mylastgame/yp-metrics-service/internal/domain/gauge"
)

type MemRepo struct {
	storage map[string]gauge.ValType
}

func NewMemRepo() *MemRepo {
	return &MemRepo{storage: make(map[string]gauge.ValType)}
}

func (r *MemRepo) Save(g *gauge.Gauge) error {
	r.storage[g.Title] = g.Val
	//fmt.Println("gauge storage: ", r.storage)
	return nil
}

func (r *MemRepo) Get(title string) (*gauge.Gauge, bool) {
	v, ok := r.storage[title]

	if ok {
		return gauge.New(title, v), ok
	} else {
		return &gauge.Gauge{}, ok
	}
}
