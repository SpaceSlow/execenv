package server

import (
	"context"
	"net/http"

	"github.com/SpaceSlow/execenv/internal/middlewares"
	"github.com/SpaceSlow/execenv/internal/routers"
	"github.com/SpaceSlow/execenv/internal/storages"
)

var _ ShutdownRunner = (*httpStrategy)(nil)

type httpStrategy struct {
	srv     *http.Server
	storage storages.MetricStorage
}

func newHttpStrategy(address string, storage storages.MetricStorage) *httpStrategy {
	runner := &httpStrategy{
		srv: &http.Server{
			Addr: address,
		},
		storage: storage,
	}
	runner.setRouters()

	return runner
}

func (s httpStrategy) Run() error {
	return s.srv.ListenAndServe()
}

func (s httpStrategy) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s httpStrategy) setRouters() {
	middlewareHandlers := []func(next http.Handler) http.Handler{
		middlewares.WithSigning,
		middlewares.WithCompressing,
		middlewares.WithDecryption,
		middlewares.WithCheckingTrustedSubnet,
		middlewares.WithLogging,
	}

	mux := routers.MetricRouter(s.storage).(http.Handler)
	for _, middleware := range middlewareHandlers {
		mux = middleware(mux)
	}

	s.srv.Handler = mux
}
