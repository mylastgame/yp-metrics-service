package handlers

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := h.repo.Ping()
	if err != nil {
		h.logger.Sugar.Error("Error connecting to database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
