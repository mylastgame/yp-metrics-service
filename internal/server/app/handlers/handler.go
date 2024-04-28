package handlers

import (
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/counter"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/gauge"
)

type Handler struct {
	GaugeRepo   gauge.Repo
	CounterRepo counter.Repo
}

func NewHandler(gr gauge.Repo, cr counter.Repo) *Handler {
	return &Handler{
		GaugeRepo:   gr,
		CounterRepo: cr,
	}
}
