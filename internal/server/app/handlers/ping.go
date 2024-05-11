package handlers

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	ps := fmt.Sprintf("host=%s user=%s password=%s",
		`localhost`, `developer`, `dev123`)

	db, err := sql.Open("pgx", ps)
	defer db.Close()

	if err != nil {
		h.logger.Sugar.Error("Error connecting to database: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
