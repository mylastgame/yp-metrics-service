package app

import (
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/collector"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/sender"
	"github.com/mylastgame/yp-metrics-service/internal/agent/storage"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
)

type app struct {
	storage   storage.Storage
	sender    sender.Sender
	collector *collector.Collector
}

func New(storage storage.Storage, sender sender.Sender, collector *collector.Collector) *app {
	return &app{
		storage:   storage,
		sender:    sender,
		collector: collector,
	}
}

func (a *app) Collect() {
	a.collector.Collect()
	logger.Log.Info("Collect finished")
}

func (a *app) Send() {
	gauges := a.storage.GetGauges()
	counters := a.storage.GetCounters()

	for t, v := range gauges {
		err := a.sender.Send(metrics.Metrics{MType: metrics.Gauge, ID: t, Value: &v})
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}

	for t, v := range counters {
		err := a.sender.Send(metrics.Metrics{MType: metrics.Counter, ID: t, Delta: &v})
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}

	logger.Log.Info("Sending finished")
	a.storage.ResetCounters()
}
