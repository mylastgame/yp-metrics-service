package handlers

import (
	"encoding/json"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
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

	if metric.MType != metrics.Counter && metric.MType != metrics.Gauge {
		logger.Log.Info("bad metric type", zap.String("type", metric.MType))
		http.Error(w, "bad metric type", http.StatusBadRequest)
	}

	err = h.repo.SaveMetric(metric)
	if err != nil {
		logger.Log.Error("Update metric error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if config.StoreInterval == 0 {
		//save data to file
		err = h.fileStorage.Save()
		if err != nil {
			logger.Log.Error("Saving to file error", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	respMetric, ok := h.repo.GetMetric(metric.MType, metric.ID)
	if !ok {
		logger.Log.Error("error when getting updated metric", zap.Error(err))
		http.Error(w, "error when getting updated metric", http.StatusInternalServerError)
		return
	}
	sendResponseMetric(w, respMetric)
}
