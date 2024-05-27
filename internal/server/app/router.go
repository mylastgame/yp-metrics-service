package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/server/app/handlers"
	"github.com/mylastgame/yp-metrics-service/internal/server/midleware"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
)

func NewRouter(repo storage.Repo, f storage.PersistentStorage, log *logger.Logger) chi.Router {
	r := chi.NewRouter()
	h := handlers.NewHandler(repo, f, log)

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

	r.Post("/update/{type}/{title}/{val}", midleware.WithLogging(midleware.GzipMiddleware(h.UpdateHandler), log))
	r.Post("/update/", midleware.WithLogging(midleware.GzipMiddleware(midleware.WithHash(h.RestUpdateHandler, log)), log))
	r.Post("/updates/", midleware.WithLogging(midleware.GzipMiddleware(midleware.WithHash(h.RestUpdatesHandler, log)), log))
	r.Get("/value/{type}/{title}", midleware.WithLogging(midleware.GzipMiddleware(midleware.WithHash(h.GetHandler, log)), log))
	r.Post("/value/", midleware.WithLogging(midleware.GzipMiddleware(midleware.WithHash(h.RestGetHandler, log)), log))
	//r.Route("/value", func(r chi.Router) {
	//	r.Get("/gauge/{title}", h.GetGaugeHandler)
	//	r.Get("/counter/{title}", h.GetCounterHandler)
	//})
	r.Get("/ping", midleware.WithLogging(midleware.GzipMiddleware(h.PingHandler), log))
	r.Get("/", midleware.WithLogging(midleware.GzipMiddleware(midleware.WithHash(h.GetAllHandler, log)), log))

	return r
}
