package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) RestUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metrics
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&metric)
	if err != nil {
		logger.Log.Error("decoding request error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if metric.MType == metrics.Gauge {
		if metric.Value == nil {
			logger.Log.Info("empty value field for gauge metric", zap.String("id", metric.ID))
			http.Error(w, fmt.Sprintf("empty value field for gauge metric: %s", metric.ID), http.StatusBadRequest)
			return
		}
		h.repo.SetGauge(metric.ID, *metric.Value)
		val, ok := h.repo.GetGauge(metric.ID)
		if !ok {
			logger.Log.Error("error when getting updated metric value", zap.Error(err))
			http.Error(w, "error when getting updated metric value", http.StatusBadRequest)
			return
		}
		sendResponseMetric(w, metrics.Metrics{ID: metric.ID, MType: metric.MType, Value: &val})
		return
	}

	if metric.MType == metrics.Counter {
		if metric.Delta == nil {
			logger.Log.Info("empty delta field for counter metric", zap.String("id", metric.ID))
			http.Error(w, fmt.Sprintf("empty delta field for counter metric: %s", metric.ID), http.StatusBadRequest)
			return
		}

		h.repo.SetCounter(metric.ID, *metric.Delta)
		val, ok := h.repo.GetCounter(metric.ID)
		if !ok {
			logger.Log.Error("error when getting updated metric value", zap.Error(err))
			http.Error(w, "error when getting updated metric value", http.StatusBadRequest)
			return
		}

		sendResponseMetric(w, metrics.Metrics{ID: metric.ID, MType: metric.MType, Delta: &val})
		return
	}

	logger.Log.Info("bad metric type", zap.String("type", metric.MType))
	http.Error(w, "bad metric type", http.StatusBadRequest)
}
