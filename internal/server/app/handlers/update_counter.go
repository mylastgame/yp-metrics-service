package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/counter"
	"io"
	"net/http"
)

func (h *Handler) UpdateCounterHandler(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	val := chi.URLParam(r, "val")
	fmt.Println("title:", title, " val:", val)

	c, err := counter.FromString(title, val)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	err = h.CounterRepo.Add(c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
