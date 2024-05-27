package handlers

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"net/http"
)

func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", config.DBConnect)
	if err != nil {
		h.logger.Sugar.Error("Error connecting to database: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			h.logger.Sugar.Error("Error closing connect to database: " + err.Error())
		}
	}()

	err = db.Ping()
	if err != nil {
		h.logger.Sugar.Error("Error connecting to database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
