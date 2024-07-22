package main

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/SpaceSlow/execenv/cmd/logger"
	"github.com/SpaceSlow/execenv/cmd/middlewares"
	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

func runServer(middlewareHandlers ...func(next http.Handler) http.Handler) error {
	if err := logger.Initialize(zap.InfoLevel.String()); err != nil {
		return err
	}

	cfg, err := GetConfigWithFlags()
	if err != nil {
		return err
	}

	var storage storages.MetricStorage
	if cfg.DatabaseDSN != "" {
		storage, err = storages.NewDBStorage(context.Background(), cfg.DatabaseDSN, cfg.Delays)
		logger.Log.Info("using storage DB", zap.String("DSN", cfg.DatabaseDSN))
	} else {
		storage, err = storages.NewMemFileStorage(cfg.StoragePath, cfg.StoreInterval, cfg.NeededRestore)
	}
	if err != nil {
		return err
	}
	defer storage.Close()
	middlewares.KEY = cfg.Key

	mux := routers.MetricRouter(storage).(http.Handler)
	for _, middleware := range middlewareHandlers {
		mux = middleware(mux)
	}

	return http.ListenAndServe(cfg.ServerAddr.String(), mux)
}

func main() {
	middlewareHandlers := []func(next http.Handler) http.Handler{
		middlewares.WithSigning,
		middlewares.WithCompressing,
		middlewares.WithLogging,
	}
	if err := runServer(middlewareHandlers...); err != nil {
		panic(err)
	}
}
