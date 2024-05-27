package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"net/http"
)

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	title := chi.URLParam(r, "title")
	ctx := r.Context()

	if t == metrics.Gauge || t == metrics.Counter {
		v, err := h.repo.Get(ctx, t, title)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(v))
			return
		}
		h.logger.Sugar.Errorf("error in Gethandler: %v", err.Error())
	}

	w.WriteHeader(http.StatusNotFound)
}
