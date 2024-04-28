package main

import (
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/collector"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/sender"
	"github.com/mylastgame/yp-metrics-service/internal/agent/config"
	"github.com/mylastgame/yp-metrics-service/internal/agent/storage"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.Create()
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)

	Sender := sender.NewHTTPSender(fmt.Sprintf("http://%s", cfg.EndpointAddr), http.MethodPost, "update")
	Storage := storage.NewMemStorage()
	App := app.New(Storage, Sender, collector.New(Storage))

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C
	sendTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	for {
		select {
		case <-pollTicker.C:
			App.Collect()

		case <-sendTicker.C:
			App.Send()
		}
	}
}
