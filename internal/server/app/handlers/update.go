package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/metrics"
	"net/http"
)

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	title := chi.URLParam(r, "title")
	val := chi.URLParam(r, "val")

	if t == metrics.Gauge || t == metrics.Counter {
		err := h.repo.Set(t, title, val)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
}
