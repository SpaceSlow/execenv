package routers

import (
	"github.com/SpaceSlow/execenv/cmd/handlers"
	"github.com/SpaceSlow/execenv/cmd/storages"
	"github.com/go-chi/chi/v5"
)

func MetricRouter(storage storages.MetricStorage) chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.BadRequestHandlerFunc)
		// router.Get("/", ...) // TODO

		r.Route("/update/", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", handlers.MetricHandler{Storage: storage}.ServeHTTP)
			r.Post("/{type}/{name}/", handlers.BadRequestHandlerFunc)
			r.Post("/", handlers.BadRequestHandlerFunc)
		})

		// router.Get("/value/{type}/{name}", ...) // TODO
	})

	return r
}
