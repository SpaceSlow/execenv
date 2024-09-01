package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/SpaceSlow/execenv/internal/config"
	"github.com/SpaceSlow/execenv/internal/logger"
	"github.com/SpaceSlow/execenv/internal/routers"
	"github.com/SpaceSlow/execenv/internal/storages"
)

func RunServer(middlewareHandlers ...func(next http.Handler) http.Handler) error {
	rootCtx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelCtx()

	if err := logger.Initialize(zap.InfoLevel.String()); err != nil {
		return err
	}

	cfg, err := config.GetServerConfig()
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(rootCtx)
	context.AfterFunc(ctx, func() {
		timeoutCtx, cancelCtx := context.WithTimeout(context.Background(), cfg.TimeoutShutdown) //nolint
		defer cancelCtx()

		<-timeoutCtx.Done()
		logger.Log.Fatal("failed to gracefully shutdown the service")
	})

	var storage storages.MetricStorage
	if cfg.DatabaseDSN != "" {
		storage, err = storages.NewDBStorage(ctx, cfg.DatabaseDSN, cfg.Delays)
		logger.Log.Info("using storage DB", zap.String("DSN", cfg.DatabaseDSN))
	} else {
		storage, err = storages.NewMemFileStorage(ctx, cfg.StoragePath, cfg.StoreInterval.Duration, cfg.NeededRestore)
	}
	if err != nil {
		return err
	}

	g.Go(func() error {
		defer logger.Log.Info("closed storage")

		<-ctx.Done()

		return storage.Close()
	})

	mux := routers.MetricRouter(storage).(http.Handler)
	for _, middleware := range middlewareHandlers {
		mux = middleware(mux)
	}
	srv := &http.Server{Addr: cfg.ServerAddr.String(), Handler: mux}

	g.Go(func() (err error) {
		defer func() {
			errRec := recover()
			if errRec != nil {
				err = fmt.Errorf("a panic occurred: %v", errRec)
			}
		}()
		if err = srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			return fmt.Errorf("listen and server has failed: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		defer logger.Log.Info("server has been shutdown")

		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), cfg.TimeoutShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			logger.Log.Error(fmt.Sprintf("an error occurred during server shutdown: %s", err))
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Log.Error(fmt.Sprintf("%s", err))
	}

	return nil
}
