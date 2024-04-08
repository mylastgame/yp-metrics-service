package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) GetCounterHandler(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	g, ok := h.CounterRepo.Get(title)

	if ok {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d", g.Val)))
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}
}
