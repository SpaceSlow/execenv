package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

var logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

func runServer() error {
	slog.SetDefault(logger)
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
