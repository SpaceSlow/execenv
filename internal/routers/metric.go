package routers

import (
	"github.com/go-chi/chi/v5"

	"github.com/SpaceSlow/execenv/internal/handlers"
	"github.com/SpaceSlow/execenv/internal/storages"
)

func MetricRouter(storage storages.MetricStorage) chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.MetricHandler{MetricStorage: storage}.List)
		r.Get("/ping", handlers.NewCheckConnectionHandler(storage).Ping)

		r.Route("/update/", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", handlers.MetricHandler{MetricStorage: storage}.Post)
			r.Post("/{type}/{name}/", handlers.BadRequestHandlerFunc)
			r.Post("/", handlers.JSONMetricHandler{MetricStorage: storage}.Post)
		})
		r.Post("/updates/", handlers.JSONMetricHandler{MetricStorage: storage}.BatchPost)
		r.Route("/value/", func(r chi.Router) {
			r.Get("/{type}/{name}", handlers.MetricHandler{MetricStorage: storage}.Get)
			r.Post("/", handlers.JSONMetricHandler{MetricStorage: storage}.Get)
		})
	})

	return r
}
