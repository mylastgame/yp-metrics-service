package main

import (
	"github.com/mylastgame/yp-metrics-service/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
