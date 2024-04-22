package main

import (
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/middlewares"
	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
	"go.uber.org/zap"
)

func runServer(middlewareHandlers ...func(next http.Handler) http.Handler) error {
	if err := middlewares.Initialize(zap.InfoLevel.String()); err != nil {
		return err
	}

	cfg, err := GetConfigWithFlags()
	if err != nil {
		return err
	}

	mux := routers.MetricRouter(storages.NewMemStorage()).(http.Handler)
	for _, middleware := range middlewareHandlers {
		mux = middleware(mux)
	}

	return http.ListenAndServe(cfg.ServerAddr.String(), mux)
}

func main() {
	middlewareHandlers := []func(next http.Handler) http.Handler{
		middlewares.WithLogging,
	}
	if err := runServer(middlewareHandlers...); err != nil {
		panic(err)
	}
}
