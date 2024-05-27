package app

import (
	"errors"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/collector"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/sender"
	"github.com/mylastgame/yp-metrics-service/internal/agent/storage"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"time"
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

func (a *app) SendBatch() {
	gauges := a.storage.GetGauges()
	counters := a.storage.GetCounters()
	metricsList := make([]metrics.Metrics, 0)

	for t, v := range gauges {
		metricsList = append(metricsList, metrics.Metrics{MType: metrics.Gauge, ID: t, Value: &v})
	}

	for t, v := range counters {
		metricsList = append(metricsList, metrics.Metrics{MType: metrics.Counter, ID: t, Delta: &v})
	}

	var errSendRequest *sender.ErrSendRequest
	tries := []int{0, 1, 3, 5}
	for i, delay := range tries {
		if i > 0 {
			time.Sleep(time.Duration(delay) * time.Second)
			a.logger.Sugar.Infof("try sending request with delay %ds", delay)
		}

		err := a.sender.SendBatch(metricsList)
		if err == nil {
			a.logger.Log.Info("sending finished")
			a.storage.ResetCounters()
			break
		}

		if errors.As(err, &errSendRequest) {
			if i < len(tries)-1 {
				a.logger.Sugar.Infof("send error: %s", err)
				continue
			}
		}

		a.logger.Sugar.Errorf("retrying number exceeded, error sending request: %s", err)
		break
	}
}
