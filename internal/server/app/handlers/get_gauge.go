package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) GetGaugeHandler(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	g, ok := h.GaugeRepo.Get(title)

	if ok {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatFloat(float64(g.Val), 'f', -1, 64)))
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}
}
