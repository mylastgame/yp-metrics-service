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
	fmt.Println("Collect finished")
}

func (a *app) Send() {
	gauges := a.storage.GetGauges()
	counters := a.storage.GetCounters()

	for t, v := range gauges {
		err := a.sender.Send(metric.Metric{Mtype: "gauge", Title: t, Val: fmt.Sprintf("%f", v)})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	for t, v := range counters {
		err := a.sender.Send(metric.Metric{Mtype: "counter", Title: t, Val: fmt.Sprintf("%d", v)})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	fmt.Println("Sending finished")

	a.storage.ResetCounters()
}
