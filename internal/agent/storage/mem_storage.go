package storage

import "sync"

type memStorage struct {
	gauges   map[string]float64
	counters map[string]int64
	m        *sync.Mutex
}

func NewMemStorage() *memStorage {
	return &memStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
		m:        &sync.Mutex{},
	}
}

func (s *memStorage) SaveGauge(title string, value float64) {
	s.m.Lock()
	defer s.m.Unlock()
	s.gauges[title] = value
}

func (s *memStorage) SaveCounter(title string, value int64) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.counters[title]; ok {
		s.counters[title] += value
	} else {
		s.counters[title] = value
	}
}

func (s *memStorage) ResetCounters() {
	s.m.Lock()
	defer s.m.Unlock()

	for k := range s.counters {
		s.counters[k] = 0
	}
}

func (s *memStorage) GetGauges() map[string]float64 {
	s.m.Lock()
	defer s.m.Unlock()

	res := make(map[string]float64)
	for k := range s.gauges {
		res[k] = s.gauges[k]
	}
	return res
}

func (s *memStorage) GetCounters() map[string]int64 {
	s.m.Lock()
	defer s.m.Unlock()

	res := make(map[string]int64)
	for k := range s.counters {
		res[k] = s.counters[k]
	}

	return res
}
