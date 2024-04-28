package server

import (
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/handlers"
	"net/http"
)

func Run() error {
	App := app.New()

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.DefaultHandler())
	mux.HandleFunc(`/update/`, handlers.UpdateHandler(App))
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		return err
	}

	return http.ListenAndServe(`:8080`, mux)
}
