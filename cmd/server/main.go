package main

import (
	"context"
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/logger"
	"github.com/SpaceSlow/execenv/cmd/middlewares"
	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func runServer(middlewareHandlers ...func(next http.Handler) http.Handler) error {
	if err := logger.Initialize(zap.InfoLevel.String()); err != nil {
		return err
	}

	cfg, err := GetConfigWithFlags()
	if err != nil {
		return err
	}

	storage, err := storages.NewMemFileStorage(cfg.StoragePath, cfg.StoreInterval, cfg.NeededRestore)
	if err != nil {
		return err
	}
	defer storage.Close()

	var conn *pgx.Conn
	if cfg.DatabaseDSN != "" {
		conn, err = pgx.Connect(context.Background(), cfg.DatabaseDSN)
		if err != nil {
			return err
		}
		defer conn.Close(context.Background())
	}

	mux := routers.MetricRouter(context.Background(), storage, conn).(http.Handler)
	for _, middleware := range middlewareHandlers {
		mux = middleware(mux)
	}

	return http.ListenAndServe(cfg.ServerAddr.String(), mux)
}

func main() {
	middlewareHandlers := []func(next http.Handler) http.Handler{
		middlewares.WithCompressing,
		middlewares.WithLogging,
	}
	if err := runServer(middlewareHandlers...); err != nil {
		panic(err)
	}
}
