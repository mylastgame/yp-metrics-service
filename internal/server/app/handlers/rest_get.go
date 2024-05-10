package handlers

import (
	"encoding/json"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) RestGetHandler(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metrics
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&metric)
	if err != nil {
		h.logger.Log.Error("decoding request error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if metric.MType != metrics.Counter && metric.MType != metrics.Gauge {
		h.logger.Log.Info("bad metric type", zap.String("type", metric.MType))
		http.Error(w, "bad metric type", http.StatusBadRequest)
	}

	if metric.ID == "" {
		h.logger.Log.Info("empty ID for gauge metric")
		http.Error(w, "empty ID for gauge metric", http.StatusNotFound)
		return
	}

	respMetric, ok := h.repo.GetMetric(metric.MType, metric.ID)
	if !ok {
		h.logger.Log.Info("metric not found", zap.String("type", metric.MType), zap.String("id", metric.ID))
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}
	h.sendResponseMetric(w, respMetric)
}
