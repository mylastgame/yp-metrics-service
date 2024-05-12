package handlers

import (
	"encoding/json"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) RestGetHandler(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metrics
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&metric)
	ctx := r.Context()

	if err != nil {
		h.logger.Log.Error("decoding request error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if metric.MType != metrics.Counter && metric.MType != metrics.Gauge {
		h.logger.Log.Info("bad metric type", zap.String("type", metric.MType))
		http.Error(w, "bad metric type", http.StatusBadRequest)
		return
	}

	if metric.ID == "" {
		h.logger.Log.Info("empty ID for gauge metric")
		http.Error(w, "empty ID for gauge metric", http.StatusNotFound)
		return
	}

	respMetric, err := h.repo.GetMetric(ctx, metric.MType, metric.ID)
	if err != nil {
		if err == storage.NotExistsError {
			h.logger.Log.Info("metric not found", zap.String("type", metric.MType), zap.String("id", metric.ID))
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}

		h.logger.Log.Error("error getting metric", zap.String("type", metric.MType), zap.String("id", metric.ID))
		http.Error(w, "metric not found", http.StatusBadRequest)
		return
	}

	h.sendResponseMetric(w, respMetric)
}
