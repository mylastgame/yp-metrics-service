package agent

import (
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/collector"
	"github.com/mylastgame/yp-metrics-service/internal/agent/app/sender"
	"github.com/mylastgame/yp-metrics-service/internal/agent/storage"
	"net/http"
	"time"
)

func Run() error {
	parseFlags()
	Sender := sender.NewHTTPSender(fmt.Sprintf("http://%s", endpointAddr), http.MethodPost, "update")
	Storage := storage.NewMemStorage()
	App := app.New(Storage, Sender, collector.New(Storage))

	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C
	sendTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)

	for {
		select {
		case <-pollTicker.C:

			App.Collect()

		case <-sendTicker.C:
			App.Send()
		}
	}
}
