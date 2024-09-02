package server

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/SpaceSlow/execenv/internal/config"
	"github.com/SpaceSlow/execenv/internal/logger"
	"github.com/SpaceSlow/execenv/internal/storages"
)

type Server struct {
	ctx context.Context

	storage        storages.MetricStorage
	config         *config.ServerConfig
	serverStrategy ShutdownRunner
}

func NewServer() (*Server, error) {
	var (
		srv Server
		err error
	)

	if err = logger.Initialize(zap.InfoLevel.String()); err != nil {
		return nil, err
	}
	srv.ctx = context.Background()

	srv.config, err = config.GetServerConfig()
	if err != nil {
		return nil, err
	}

	err = srv.setStorage()
	if err != nil {
		return nil, err
	}

	srv.setStrategy()

	return &srv, nil
}

func (s *Server) Run() error {
	rootCtx, cancelCtx := signal.NotifyContext(s.ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelCtx()

	g, ctx := errgroup.WithContext(rootCtx)
	s.ctx = ctx

	context.AfterFunc(ctx, func() {
		timeoutCtx, cancelCtx := context.WithTimeout(context.Background(), s.config.TimeoutShutdown)
		defer cancelCtx()

		<-timeoutCtx.Done()
		logger.Log.Fatal("failed to gracefully shutdown the service")
	})

	g.Go(func() error {
		defer logger.Log.Info("closed storage")

		<-ctx.Done()

		return s.storage.Close()
	})

	g.Go(func() (err error) {
		defer func() {
			errRec := recover()
			if errRec != nil {
				err = fmt.Errorf("a panic occurred: %v", errRec)
			}
		}()
		return s.serverStrategy.Run()
	})

	g.Go(func() error {
		defer logger.Log.Info("server has been shutdown")

		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), s.config.TimeoutShutdown)
		defer cancelShutdownTimeoutCtx()
		return s.serverStrategy.Shutdown(shutdownTimeoutCtx)
	})

	if err := g.Wait(); err != nil {
		logger.Log.Error(fmt.Sprintf("%s", err))
	}

	return nil
}

func (s *Server) setStorage() error {
	var err error

	if s.config.DatabaseDSN != "" {
		s.storage, err = storages.NewDBStorage(s.ctx, s.config.DatabaseDSN, s.config.Delays)
		logger.Log.Info("using storage DB", zap.String("DSN", s.config.DatabaseDSN))
	} else {
		s.storage, err = storages.NewMemFileStorage(s.ctx, s.config.StoragePath, s.config.StoreInterval.Duration, s.config.NeededRestore)
	}

	return err
}

func (s *Server) setStrategy() {
	if s.config.StartedGRPCServer {
		s.serverStrategy = newGrpcStrategy(s.config.ServerAddr.String(), s.storage)
		return
	}
	s.serverStrategy = newHttpStrategy(s.config.ServerAddr.String(), s.storage)
}
