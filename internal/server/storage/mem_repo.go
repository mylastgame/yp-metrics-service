package storage

import (
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/service/convert"
	"sync"
)

type MemRepo struct {
	gauge   map[string]float64
	counter map[string]int64
	m       *sync.Mutex
}

func NewMemRepo() *MemRepo {
	return &MemRepo{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
		m:       &sync.Mutex{},
	}
}

func (r *MemRepo) Get(t string, k string) (string, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if t == metrics.Gauge {
		v, ok := r.gauge[k]
		if ok {
			return convert.GaugeToString(v), nil
		} else {
			return "", NewStorageError(KeyNotExists, t, k)
		}
	}

	if t == metrics.Counter {
		v, ok := r.counter[k]
		if ok {
			return convert.CounterToString(v), nil
		} else {
			return "", NewStorageError(KeyNotExists, t, k)
		}
	}

	return "", NewStorageError(BadMetricType, t, k)
}

func (r *MemRepo) Set(t string, k string, v string) error {
	r.m.Lock()
	defer r.m.Unlock()

	if t == metrics.Gauge {
		g, err := convert.StringToGauge(v)
		if err != nil {
			return NewStorageError(BadValue, t, v)
		}
		r.gauge[k] = g
		return nil
	}

	if t == metrics.Counter {
		c, err := convert.StringToCounter(v)
		if err != nil {
			return NewStorageError(BadValue, t, v)
		}

		_, ok := r.counter[k]

		if ok {
			r.counter[k] += c
		} else {
			r.counter[k] = c
		}

		return nil
	}

	return NewStorageError(BadMetricType, t, k)
}

func (r *MemRepo) GetCounters() []string {
	r.m.Lock()
	defer r.m.Unlock()

	res := make([]string, 0)
	for _, v := range r.counter {
		res = append(res, convert.CounterToString(v))
	}

	return res
}

func (r *MemRepo) GetGauges() []string {
	r.m.Lock()
	defer r.m.Unlock()

	res := make([]string, 0)
	for _, v := range r.gauge {
		res = append(res, convert.GaugeToString(v))
	}

	return res
}
