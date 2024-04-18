package main

import (
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/middlewares"
	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
	"go.uber.org/zap"
)

func runServer() error {
	if err := middlewares.Initialize(zap.InfoLevel.String()); err != nil {
		return err
	}

	cfg, err := GetConfigWithFlags()
	if err != nil {
		return err
	}
	return http.ListenAndServe(cfg.ServerAddr.String(), routers.MetricRouter(storages.NewMemStorage()))
}

func main() {
	if err := runServer(); err != nil {
		panic(err)
	}
}
