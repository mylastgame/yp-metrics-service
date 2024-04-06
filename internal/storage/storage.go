package storage

import metrics "github.com/mylastgame/yp-metrics-service/internal/metrics"

type StorageI interface {
	AddCounter(string, metrics.Counter)
	AddGauge(string, metrics.Gauge)
}
