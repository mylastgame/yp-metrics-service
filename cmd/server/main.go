package main

import (
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"net/http"
)

func main() {
	config.ParseFlags()
	err := logger.Initialize(config.LogLevel)
	if err != nil {
		panic(err)
	}

	r := app.NewRouter(storage.NewMemRepo())

	logger.Sugar.Infof("Starting server. Listening on %s", config.RunAddr)
	err = http.ListenAndServe(config.RunAddr, r)

	if err != nil {
		panic(err)
	}
}
