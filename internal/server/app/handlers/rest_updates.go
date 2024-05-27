package handlers

import (
	"encoding/json"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) RestUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	var metrics []metrics.Metrics
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&metrics)
	ctx := r.Context()

	if err != nil {
		h.logger.Log.Error("decoding request error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repo.SaveMetrics(ctx, metrics)
	if err != nil {
		h.logger.Log.Error("saving metrics error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Sugar.Infof("metrics updated: %v", metrics)
	w.WriteHeader(http.StatusOK)
}
