package handlers

import (
	"encoding/json"
	"errors"
	"github.com/mylastgame/yp-metrics-service/internal/core"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) RestUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	err = core.Retry("save metric", 3, func() error {
		return h.repo.SaveMetric(ctx, metric)
	}, h.logger)

	if err != nil {
		h.logger.Log.Error("Update metric error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if config.StoreInterval == 0 {
		//save data to file
		err = h.fileStorage.Save(ctx)
		if err != nil {
			h.logger.Log.Error("Saving to file error", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var respMetric metrics.Metrics
	err = core.Retry("get response metric", 3, func() error {
		respMetric, err = h.repo.GetMetric(ctx, metric.MType, metric.ID)
		return err
	}, h.logger)
	if err != nil {
		if errors.Is(err, storage.ErrorNotExists) {
			h.logger.Log.Error("error when getting updated metric", zap.Error(err))
			http.Error(w, "error when getting updated metric", http.StatusInternalServerError)
			return
		}

		h.logger.Log.Error("error getting metric", zap.String("type", metric.MType), zap.String("id", metric.ID))
		http.Error(w, "metric not found", http.StatusBadRequest)
		return
	}
	h.sendResponseMetric(w, respMetric)
}
