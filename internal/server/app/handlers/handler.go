package handlers

import (
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
)

type Handler struct {
	repo storage.Repo
}

func NewHandler(r storage.Repo) *Handler {
	return &Handler{repo: r}
}
