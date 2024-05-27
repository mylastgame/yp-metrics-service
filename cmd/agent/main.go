package main

import (
	"context"
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/collector"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/sender"
	"github.com/mylastgame/yp-metrics-service/internal/agent/config"
	"github.com/mylastgame/yp-metrics-service/internal/agent/storage"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.Create()
	if err != nil {
		panic(err)
	}

	//logger init
	log, err := logger.NewLogger("info")
	if err != nil {
		fmt.Printf("Error init logger: %v/n", err)
		panic(err)
	}

	//Sender := sender.NewHTTPSender(fmt.Sprintf("http://%s", cfg.EndpointAddr), http.MethodPost, "update")
	Sender := sender.NewRESTSender(fmt.Sprintf("http://%s", cfg.EndpointAddr), http.MethodPost, "updates", log, &cfg)
	Storage := storage.NewMemStorage()
	App := app.New(Storage, Sender, collector.New(Storage), log)
	log.Log.Sugar().Infof("Agent started. Poll interval: %ds, report interval: %ds, endpoint: %s, key: %s",
		cfg.PollInterval, cfg.ReportInterval, cfg.EndpointAddr, cfg.Key)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C
	sendTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	for {
		select {
		case <-pollTicker.C:
			App.Collect()
		case <-sendTicker.C:
			//	App.Send()
			App.SendBatch()
		case <-ctx.Done():
			return
		}
	}
}
