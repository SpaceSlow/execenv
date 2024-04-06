package main

import (
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

func runServer() error {
	return http.ListenAndServe("localhost:8080", routers.MetricRouter(storages.NewMemStorage()))
}

func main() {
	if err := runServer(); err != nil {
		panic(err)
	}
}
