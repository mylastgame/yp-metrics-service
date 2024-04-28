package handlers

import (
	"github.com/mylastgame/yp-metrics-service/internal/service/html"
	"net/http"
)

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	gauges := h.repo.GetGauges()
	counters := h.repo.GetCounters()

	gaugeHTML := "Gauges: <ol>"
	//html.SliceToOlLi("Gauges", gauges)
	for _, g := range gauges {
		gaugeHTML += html.Tag("li", g)
	}
	gaugeHTML += "</ol>"

	counterHTML := "Counters: <ol>"
	for _, c := range counters {
		counterHTML += html.Tag("li", c)
	}
	counterHTML += "</ol>"

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(gaugeHTML + counterHTML))
}
