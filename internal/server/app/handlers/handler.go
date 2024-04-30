package handlers

import (
	"encoding/json"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	repo storage.Repo
}

func NewHandler(r storage.Repo) *Handler {
	return &Handler{repo: r}
}

func sendResponseMetric(w http.ResponseWriter, metric metrics.Metrics) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	err := enc.Encode(metric)
	if err != nil {
		logger.Log.Error("encoding response error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var field zap.Field
	if metric.MType == metrics.Gauge {
		field = zap.Float64("value", *metric.Value)
	} else if metric.MType == metrics.Counter {
		field = zap.Int64("delta", *metric.Delta)
	}

	logger.Log.Info("metric updated",
		zap.String("type", metric.MType),
		zap.String("id", metric.ID),
		field,
	)
}
