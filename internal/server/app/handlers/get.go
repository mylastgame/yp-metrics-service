package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/metrics"
	"net/http"
)

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	title := chi.URLParam(r, "title")

	if t == metrics.Gauge || t == metrics.Counter {
		v, err := h.repo.Get(t, title)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(v))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
