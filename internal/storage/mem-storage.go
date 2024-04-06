package storage

import "github.com/mylastgame/yp-metrics-service/internal/metrics"

type MemStorage struct {
	Counters map[string]metrics.Counter
	Gauges   map[string]metrics.Gauge
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Counters: make(map[string]metrics.Counter),
		Gauges:   make(map[string]metrics.Gauge),
	}
}

func (ms *MemStorage) AddCounter(key string, value metrics.Counter) {
	if _, ok := ms.Counters[key]; ok {
		ms.Counters[key] += value
	} else {
		ms.Counters[key] = value
	}
}

func (ms *MemStorage) AddGauge(key string, value metrics.Gauge) {
	ms.Gauges[key] = value
}
