package main

import (
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/flags"
	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

func runServer(addr flags.NetAddress) error {
	return http.ListenAndServe(addr.String(), routers.MetricRouter(storages.NewMemStorage()))
}

func main() {
	parseFlags()

	if err := runServer(flagRunAddr); err != nil {
		panic(err)
	}
}
