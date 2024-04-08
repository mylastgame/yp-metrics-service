package handlers

import (
	"github.com/mylastgame/yp-metrics-service/internal/service/html"
	"net/http"
)

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	gauges := h.GaugeRepo.GetAll()
	counters := h.CounterRepo.GetAll()

	gaugeHtml := "Gauges: <ol>"
	//html.SliceToOlLi("Gauges", gauges)
	for _, g := range gauges {
		gaugeHtml += html.Tag("li", g)
	}
	gaugeHtml += "</ol>"

	counterHtml := "Counters: <ol>"
	for _, c := range counters {
		counterHtml += html.Tag("li", c)
	}
	counterHtml += "</ol>"

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(gaugeHtml + counterHtml))
	return
}
