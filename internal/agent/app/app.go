package app

import (
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/collector"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/sender"
	"github.com/mylastgame/yp-metrics-service/internal/agent/metric"
	"github.com/mylastgame/yp-metrics-service/internal/agent/storage"
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
}

func (a *app) Send() {
	gauges := a.storage.GetGauges()
	counters := a.storage.GetCounters()

	for t, v := range gauges {
		err := a.sender.Send(metric.Metric{"gauge", t, fmt.Sprintf("%f", v)})
		if err != nil {
			panic(err)
		}
	}

	for t, v := range counters {
		err := a.sender.Send(metric.Metric{"counter", t, fmt.Sprintf("%d", v)})
		if err != nil {
			panic(err)
		}
	}

	a.storage.ResetCounters()

	fmt.Printf("Send\n")
}
