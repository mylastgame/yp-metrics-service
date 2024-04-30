package handlers

import (
	"encoding/json"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) RestGetHandler(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metrics
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&metric)
	if err != nil {
		logger.Log.Error("decoding request error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if metric.MType == metrics.Gauge {
		if metric.ID == "" {
			logger.Log.Info("empty ID for gauge metric")
			http.Error(w, "empty ID for gauge metric", http.StatusNotFound)
			return
		}
		val, ok := h.repo.GetGauge(metric.ID)
		if !ok {
			logger.Log.Info("metric not found", zap.String("type", metric.MType), zap.String("id", metric.ID))
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		sendResponseMetric(w, metrics.Metrics{ID: metric.ID, MType: metric.MType, Value: &val})
		return
	}

	if metric.MType == metrics.Counter {
		if metric.ID == "" {
			logger.Log.Info("empty ID for counter metric")
			http.Error(w, "empty ID for counter metric", http.StatusNotFound)
			return
		}
		var val float64
		v, ok := h.repo.GetCounter(metric.ID)
		if !ok {
			logger.Log.Info("metric not found", zap.String("type", metric.MType), zap.String("id", metric.ID))
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		} else {
			val = float64(v)
		}
		sendResponseMetric(w, metrics.Metrics{ID: metric.ID, MType: metric.MType, Value: &val})
		return
	}

	logger.Log.Info("bad metric type", zap.String("type", metric.MType))
	http.Error(w, "bad metric type", http.StatusBadRequest)
}
