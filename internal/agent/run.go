package agent

import (
	"github.com/mylastgame/yp-metrics-service/internal/agent/app"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/collector"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/sender"
	"github.com/mylastgame/yp-metrics-service/internal/agent/storage"
	"net/http"
	"time"
)

func Run() error {

	Sender := sender.NewHTTPSender("http://localhost:8080", http.MethodPost, "update")
	Storage := storage.NewMemStorage()
	App := app.New(Storage, Sender, collector.New(Storage))

	pollTicker := time.NewTicker(2 * time.Second)
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C
	sendTicker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-pollTicker.C:

			App.Collect()

		case <-sendTicker.C:
			App.Send()
		}
	}
}
