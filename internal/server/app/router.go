package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/server/app/handlers"
	"github.com/mylastgame/yp-metrics-service/internal/server/midleware"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
)

func NewRouter(repo storage.Repo) chi.Router {
	r := chi.NewRouter()
	h := handlers.NewHandler(repo)

	//r.Route("/update", func(r chi.Router) {
	//	r.Route("/gauge", func(r chi.Router) {
	//		r.Post("/", handlers.NotFoundHandler())
	//		r.Post("/{title}/{val}", h.UpdateGaugeHandler)
	//	})
	//	r.Route("/counter", func(r chi.Router) {
	//		r.Post("/", handlers.NotFoundHandler())
	//		r.Post("/{title}/{val}", h.UpdateCounterHandler)
	//	})
	//	r.Post("/{type}/{title}/{val}", handlers.BadRequestHandler())
	//})

	r.Post("/update/{type}/{title}/{val}", midleware.WithLogging(midleware.GzipMiddleware(h.UpdateHandler)))
	r.Post("/update/", midleware.WithLogging(midleware.GzipMiddleware(h.RestUpdateHandler)))
	r.Get("/value/{type}/{title}", midleware.WithLogging(midleware.GzipMiddleware(h.GetHandler)))
	r.Post("/value/", midleware.WithLogging(midleware.GzipMiddleware(h.RestGetHandler)))
	//r.Route("/value", func(r chi.Router) {
	//	r.Get("/gauge/{title}", h.GetGaugeHandler)
	//	r.Get("/counter/{title}", h.GetCounterHandler)
	//})

	r.Get("/", midleware.WithLogging(midleware.GzipMiddleware(h.GetAllHandler)))

	return r
}
