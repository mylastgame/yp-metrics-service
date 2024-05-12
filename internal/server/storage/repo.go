package storage

import (
	"context"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
)

type Repo interface {
	Set(context.Context, string, string, string) error
	Get(context.Context, string, string) (string, error)
	SetGauge(context.Context, string, float64) error
	SetCounter(context.Context, string, int64) error
	GetGauge(context.Context, string) (float64, error)
	GetCounter(context.Context, string) (int64, error)
	GetGauges(context.Context) (metrics.GaugeList, error)
	GetCounters(context.Context) (metrics.CounterList, error)
	SaveMetric(context.Context, metrics.Metrics) error
	GetMetric(context.Context, string, string) (metrics.Metrics, error)
	Ping() error
}
