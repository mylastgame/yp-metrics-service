package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/server/app/handlers"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/counter"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/gauge"
)

func Setup(r *chi.Mux, gr gauge.Repo, cr counter.Repo) {
	h := handlers.NewHandler(gr, cr)

	r.Route("/update", func(r chi.Router) {
		r.Route("/gauge", func(r chi.Router) {
			r.Post("/", handlers.NotFoundHandler())
			r.Post("/{title}/{val}", h.UpdateGaugeHandler)
		})
		r.Route("/counter", func(r chi.Router) {
			r.Post("/", handlers.NotFoundHandler())
			r.Post("/{title}/{val}", h.UpdateCounterHandler)
		})
		r.Post("/{type}/{title}/{val}", handlers.BadRequestHandler())
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/gauge/{title}", h.GetGaugeHandler)
		r.Get("/counter/{title}", h.GetCounterHandler)
	})

	r.Get("/", h.GetHandler)
}
