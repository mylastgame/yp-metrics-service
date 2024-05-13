package storage

import (
	"context"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
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

func (r *MemRepo) Get(ctx context.Context, t string, k string) (string, error) {
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

func (r *MemRepo) Set(ctx context.Context, t string, k string, v string) error {
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

func (r *MemRepo) SetGauge(ctx context.Context, k string, v float64) error {
	r.m.Lock()
	defer r.m.Unlock()
	r.gauge[k] = v
	return nil
}

func (r *MemRepo) SetCounter(ctx context.Context, k string, v int64) error {
	r.m.Lock()
	defer r.m.Unlock()

	r.counter[k] += v
	//_, ok := r.counter[k]
	//
	//if ok {
	//	r.counter[k] += v
	//} else {
	//	r.counter[k] = v
	//}

	return nil
}

func (r *MemRepo) GetCounter(ctx context.Context, k string) (int64, error) {
	r.m.Lock()
	defer r.m.Unlock()

	v, ok := r.counter[k]

	if ok {
		return v, nil
	} else {
		return 0, NewStorageError(KeyNotExists, metrics.Counter, k)
	}
}

func (r *MemRepo) GetGauge(ctx context.Context, k string) (float64, error) {
	r.m.Lock()
	defer r.m.Unlock()

	v, ok := r.gauge[k]

	if ok {
		return v, nil
	} else {
		return 0, NewStorageError(KeyNotExists, metrics.Gauge, k)
	}
}

func (r *MemRepo) GetGauges(ctx context.Context) (metrics.GaugeList, error) {
	r.m.Lock()
	defer r.m.Unlock()

	res := metrics.GaugeList{}
	for k, v := range r.gauge {
		res[k] = v
	}

	return res, nil
}

func (r *MemRepo) GetCounters(ctx context.Context) (metrics.CounterList, error) {
	r.m.Lock()
	defer r.m.Unlock()

	res := metrics.CounterList{}
	for k, v := range r.counter {
		res[k] = v
	}

	return res, nil
}

func (r *MemRepo) SaveMetric(ctx context.Context, metric metrics.Metrics) error {
	if metric.MType == metrics.Gauge {
		err := r.SetGauge(ctx, metric.ID, *metric.Value)
		if err != nil {
			return err
		}
		return nil
	}

	if metric.MType == metrics.Counter {
		err := r.SetCounter(ctx, metric.ID, *metric.Delta)
		if err != nil {
			return err
		}
		return nil
	}

	return NewStorageError(BadMetricType, metric.MType, metric.ID)
}

func (r *MemRepo) GetMetric(ctx context.Context, mType string, id string) (metrics.Metrics, error) {
	metric := metrics.Metrics{}
	var (
		err  error
		gVal float64
		cVal int64
	)

	if mType == metrics.Gauge {
		gVal, err = r.GetGauge(ctx, id)
		if err == nil {
			metric.MType = mType
			metric.ID = id
			metric.Value = &gVal
		}
		return metric, err
	}

	if mType == metrics.Counter {
		cVal, err = r.GetCounter(ctx, id)
		if err == nil {
			metric.MType = mType
			metric.ID = id
			metric.Delta = &cVal
		}
		return metric, err
	}

	return metric, NewStorageError(BadMetricType, metric.MType, metric.ID)
}

func (r *MemRepo) SaveMetrics(ctx context.Context, list []metrics.Metrics) error {
	for _, metric := range list {
		err := r.SaveMetric(ctx, metric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MemRepo) Ping() error {
	return nil
}
