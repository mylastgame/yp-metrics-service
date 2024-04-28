package main

import "github.com/mylastgame/yp-metrics-service/internal/agent"

func main() {
	err := agent.Run()
	if err != nil {
		panic(err)
	}

}
