package storage

import "github.com/mylastgame/yp-metrics-service/internal/core/metrics"

type Repo interface {
	Set(string, string, string) error
	Get(string, string) (string, error)
	SetGauge(string, float64)
	SetCounter(string, int64)
	GetGauge(string) (float64, bool)
	GetCounter(string) (int64, bool)
	GetGauges() map[string]float64
	GetCounters() map[string]int64
	SaveMetric(metrics.Metrics) error
	GetMetric(string, string) (metrics.Metrics, bool)
}
