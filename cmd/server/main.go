package main

import (
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

func runServer() error {
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
