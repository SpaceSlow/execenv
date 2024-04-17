package routers

import (
	"github.com/SpaceSlow/execenv/cmd/handlers"
	"github.com/SpaceSlow/execenv/cmd/middlewares"
	"github.com/SpaceSlow/execenv/cmd/storages"
	"github.com/go-chi/chi/v5"
)

func MetricRouter(storage storages.MetricStorage) chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", middlewares.WithLogging(handlers.BadRequestHandlerFunc))
		r.Get("/", middlewares.WithLogging(handlers.MetricHandler{MetricStorage: storage}.List))

		r.Route("/update/", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", middlewares.WithLogging(handlers.MetricHandler{MetricStorage: storage}.Post))
			r.Post("/{type}/{name}/", middlewares.WithLogging(handlers.BadRequestHandlerFunc))
			r.Post("/", middlewares.WithLogging(handlers.BadRequestHandlerFunc))
		})

		r.Get("/value/{type}/{name}", middlewares.WithLogging(handlers.MetricHandler{MetricStorage: storage}.Get))
	})

	return r
}
