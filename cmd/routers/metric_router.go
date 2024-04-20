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
		r.Get("/", middlewares.WithLogging(handlers.MetricHandler{MetricStorage: storage}.List))

		r.Route("/update/", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", middlewares.WithLogging(handlers.MetricHandler{MetricStorage: storage}.Post))
			r.Post("/{type}/{name}/", middlewares.WithLogging(handlers.BadRequestHandlerFunc))
			r.Post("/", middlewares.WithLogging(handlers.JSONMetricHandler{MetricStorage: storage}.Post))
		})
		r.Route("/value/", func(r chi.Router) {
			r.Get("/{type}/{name}", middlewares.WithLogging(handlers.MetricHandler{MetricStorage: storage}.Get))
			r.Post("/", middlewares.WithLogging(handlers.JSONMetricHandler{MetricStorage: storage}.Get))
		})
	})

	return r
}
