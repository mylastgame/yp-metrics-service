package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage/counter"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage/gauge"
	"net/http"
)

func Run() error {
	parseFlags()

	r := chi.NewRouter()
	app.Setup(r, gauge.NewMemRepo(), counter.NewMemRepo())

	fmt.Printf("Listening on %s\n", flagRunAddr)

	err := http.ListenAndServe(flagRunAddr, r)
	if err != nil {
		return err
	}

	return nil
}
