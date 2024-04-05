package main

import (
	"fmt"
	metrics "github.com/mylastgame/yp-metrics-service/internal/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/storage"
	"io"
	"net/http"
	"strings"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, defaultPage)
	mux.HandleFunc(`/update/`, updateMetrics)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

	return http.ListenAndServe(`:8080`, mux)
}

func defaultPage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "YPrac Metrics and alerting service")
}

func updateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlPath := strings.Split(r.URL.Path, "/")
	urlLen := len(urlPath)

	//if urlLen < 5 {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	//check type exists
	if !metrics.TypeExists(urlPath[2]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//check title isset
	if urlLen == 3 || urlPath[3] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//check title exists
	if !metrics.MetricExists(urlPath[2], urlPath[3]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//check value isset
	if urlLen == 4 || urlPath[4] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r, _ := metrics.Save(urlPath[2], urlPath[3], urlPath[4]); r {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, fmt.Sprintf("storage: %v\n", storage.Storage))
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
