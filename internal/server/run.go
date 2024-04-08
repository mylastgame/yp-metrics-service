package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage/counter"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage/gauge"
	"net/http"
)

func Run() error {
	r := chi.NewRouter()
	app.Setup(r, gauge.NewMemRepo(), counter.NewMemRepo())

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		return err
	}

	return nil
}
