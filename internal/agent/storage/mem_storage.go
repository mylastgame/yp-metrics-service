package storage

type memStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *memStorage {
	return &memStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *memStorage) SaveGauge(title string, value float64) {
	s.gauges[title] = value
}

func (s *memStorage) SaveCounter(title string, value int64) {
	_, ok := s.counters[title]

	if ok {
		s.counters[title] += value
	} else {
		s.counters[title] = value
	}
}

func (s *memStorage) ResetCounters() {
	for k, _ := range s.counters {
		s.counters[k] = 0
	}
}

func (s *memStorage) GetGauges() map[string]float64 {
	res := make(map[string]float64)
	for k, _ := range s.gauges {
		res[k] = s.gauges[k]
	}
	return res
}

func (s *memStorage) GetCounters() map[string]int64 {
	res := make(map[string]int64)
	for k, _ := range s.counters {
		res[k] = s.counters[k]
	}

	return res
}
