package app

import (
	"context"
	"errors"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/collector"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/sender"
	"github.com/mylastgame/yp-metrics-service/internal/agent/config"
	"github.com/mylastgame/yp-metrics-service/internal/agent/storage"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"sync"
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

func (a *app) CollectGOPSUtil() {
	err := a.collector.CollectGOPSUtil()
	if err != nil {
		a.logger.Log.Error(err.Error())
		return
	}
	a.logger.Log.Info("Collect GOPSUtil finished")
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

func (a *app) Run(ctx context.Context, cfg *config.Config) error {
	wg := &sync.WaitGroup{}

	//collecting
	pollTicker1 := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollTicker1.Stop()
	a.runJob(ctx, wg, pollTicker1.C, func() { a.Collect() })

	//collecting GOPSUtil
	pollTicker2 := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollTicker2.Stop()
	a.runJob(ctx, wg, pollTicker2.C, func() { a.CollectGOPSUtil() })

	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C

	sendTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer sendTicker.Stop()

	// создаем буферизованный канал для принятия задач
	jobsCh := make(chan time.Time, cfg.RateLimit*3)
	defer close(jobsCh)

	//create worker pool
	a.startWorkerPool(ctx, *cfg, jobsCh)
	//sending data to server
	a.runJob(ctx, wg, sendTicker.C, func() { jobsCh <- time.Now() })

	wg.Wait()
	return nil
}

func (a *app) runJob(ctx context.Context, wg *sync.WaitGroup, c <-chan time.Time, job func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case _, ok := <-c:
				if !ok {
					return
				}
				job()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// create worker pool
func (a *app) startWorkerPool(ctx context.Context, cfg config.Config, jobs <-chan time.Time) {
	//create working pool
	for w := 0; w < cfg.RateLimit; w++ {
		go func() {
			for {
				select {
				case _, ok := <-jobs:
					if !ok {
						return
					}
					a.SendBatch()
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	a.logger.Sugar.Infof("Worker pool started with limit: %d", cfg.RateLimit)
}
