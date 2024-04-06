package app

import (
	"github.com/mylastgame/yp-metrics-service/internal/server/storagesrv"
	"github.com/mylastgame/yp-metrics-service/internal/storage"
	"io"
	"net/http"
	"strings"
)

type App struct {
	srv storagesrv.StorageServiceI
}

func New() *App {
	return &App{srv: storagesrv.New(storage.NewMemStorage())}
}

func (app *App) Run() error {
	defaultPage := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updateMetrics := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			// разрешаем только POST-запросы
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		urlPath := strings.Split(r.URL.Path, "/")
		urlPathLen := len(urlPath)

		//if urlPathLen != 5 {
		//	w.WriteHeader(http.StatusBadRequest)
		//	return
		//}

		if urlPathLen < 4 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mtype := urlPath[2]
		mtitle := urlPath[3]

		//check type in path
		if mtype == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//check title in path
		if mtitle == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if urlPathLen < 5 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mval := urlPath[4]

		if mval == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := app.srv.Save(mtype, mtitle, mval); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, defaultPage)
	mux.HandleFunc(`/update/`, updateMetrics)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		return err
	}

	return http.ListenAndServe(`:8080`, mux)
}
