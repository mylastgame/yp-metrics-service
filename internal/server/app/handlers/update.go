package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	title := chi.URLParam(r, "title")
	val := chi.URLParam(r, "val")

	if t == metrics.Gauge || t == metrics.Counter {
		err := h.repo.Set(t, title, val)
		if err == nil {
			if config.StoreInterval == 0 {
				//save data to file
				err = h.fileStorage.Save()
				if err != nil {
					h.logger.Log.Error("Saving to file error", zap.Error(err))
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
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
