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
	logger    *logger.Logger
}

func New(storage storage.Storage, sender sender.Sender, collector *collector.Collector, log *logger.Logger) *app {
	return &app{
		storage:   storage,
		sender:    sender,
		collector: collector,
		logger:    log,
	}
}

func (a *app) Collect() {
	a.collector.Collect()
	a.logger.Log.Info("Collect finished")
}

func (a *app) Send() {
	gauges := a.storage.GetGauges()
	counters := a.storage.GetCounters()

	for t, v := range gauges {
		err := a.sender.Send(metrics.Metrics{MType: metrics.Gauge, ID: t, Value: &v})
		if err != nil {
			a.logger.Log.Error(err.Error())
		}
	}

	for t, v := range counters {
		err := a.sender.Send(metrics.Metrics{MType: metrics.Counter, ID: t, Delta: &v})
		if err != nil {
			a.logger.Log.Error(err.Error())
		}
	}

	a.logger.Log.Info("Sending finished")
	a.storage.ResetCounters()
}
