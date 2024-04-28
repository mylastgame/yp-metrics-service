package main

import (
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"net/http"
)

func main() {
	config.ParseFlags()
	fmt.Println(config.RunAddr)

	r := app.NewRouter(storage.NewMemRepo())

	err := http.ListenAndServe(config.RunAddr, r)
	if err != nil {
		panic(err)
	}
}
