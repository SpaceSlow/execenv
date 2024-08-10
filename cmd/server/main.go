package main

import (
	"context"
	"go.uber.org/zap"
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/config"
	"github.com/SpaceSlow/execenv/cmd/logger"
	"github.com/SpaceSlow/execenv/cmd/middlewares"
	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

func RunServer(middlewareHandlers ...func(next http.Handler) http.Handler) error {
	if err := logger.Initialize(zap.InfoLevel.String()); err != nil {
		return err
	}

	cfg, err := config.GetServerConfig()
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

	mux := routers.MetricRouter(storage).(http.Handler)
	for _, middleware := range middlewareHandlers {
		mux = middleware(mux)
	}

	return http.ListenAndServe(cfg.ServerAddr.String(), mux)
}

func main() {
	config.PrintBuildInfo()
	middlewareHandlers := []func(next http.Handler) http.Handler{
		middlewares.WithSigning,
		middlewares.WithCompressing,
		middlewares.WithDecryption,
		middlewares.WithLogging,
	}
	if err := RunServer(middlewareHandlers...); err != nil {
		panic(err)
	}
}
