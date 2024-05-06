package handlers

import (
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/service/html"
	"net/http"
)

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	gauges := h.repo.GetGauges()
	counters := h.repo.GetCounters()

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
