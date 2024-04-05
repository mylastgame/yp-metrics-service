package storage

type MemStorage struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

var Storage *MemStorage

func init() {
	Storage = &MemStorage{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}
}

func (ms *MemStorage) AddCounter(key string, value int64) {
	if _, ok := ms.Counters[key]; ok {
		ms.Counters[key] += value
	} else {
		ms.Counters[key] = value
	}
}

func (ms *MemStorage) AddGauge(key string, value float64) {
	ms.Gauges[key] = value
}
