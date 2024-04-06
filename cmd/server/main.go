package main

import (
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
)

func main() {
	App := app.New()
	if err := App.Run(); err != nil {
		panic(err)
	}
}
