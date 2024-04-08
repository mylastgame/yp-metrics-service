package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/gauge"
	"io"
	"net/http"
)

func (h *Handler) UpdateGaugeHandler(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	val := chi.URLParam(r, "val")
	fmt.Println("title:", title, " val:", val)

	g, err := gauge.FromString(title, val)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	err = h.GaugeRepo.Save(g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
