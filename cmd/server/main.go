package main

import (
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/handlers"
	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

func runServer() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.DefaultHandler)
	mux.Handle("/update/counter/", http.StripPrefix("/update/counter/", handlers.MetricHandler{
		MetricType: metrics.Counter,
		Storage:    storages.NewMemStorage(),
	}))
	mux.Handle("/update/gauge/", http.StripPrefix("/update/gauge/", handlers.MetricHandler{
		MetricType: metrics.Gauge,
		Storage:    storages.NewMemStorage(),
	}))

	return http.ListenAndServe("localhost:8080", mux)
}

func main() {
	if err := runServer(); err != nil {
		panic(err)
	}
}
