package handlers

import (
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"github.com/mylastgame/yp-metrics-service/internal/service/html"
	"net/http"
)

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var (
		err      error
		gauges   metrics.GaugeList
		counters metrics.CounterList
	)

	gauges, err = h.repo.GetGauges(ctx)
	if err != nil && err != storage.NotExistsError {
		h.logger.Sugar.Errorf("GetAllHandler: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	counters, err = h.repo.GetCounters(ctx)
	if err != nil && err != storage.NotExistsError {
		h.logger.Sugar.Errorf("GetAllHandler: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	gaugeHTML := "Gauges: <ol>"
	//html.SliceToOlLi("Gauges", gauges)
	for k, g := range gauges {
		gaugeHTML += html.Tag("li", fmt.Sprintf("%s: %f", k, g))
	}
	gaugeHTML += "</ol>"

	counterHTML := "Counters: <ol>"
	for k, c := range counters {
		counterHTML += html.Tag("li", fmt.Sprintf("%s: %d", k, c))
	}
	counterHTML += "</ol>"

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(gaugeHTML + counterHTML))
}
