package main

import (
	"github.com/SpaceSlow/execenv/cmd/server/metrics"
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/server/handlers"
	"github.com/SpaceSlow/execenv/cmd/server/storages"
)

func runServer() error {
	storage := storages.NewMemStorage()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.DefaultHandler)
	mux.Handle("/update/counter/", http.StripPrefix("/update/counter/", handlers.MetricHandler{
		MetricType: metrics.Counter,
		Storage:    storage,
	}))
	mux.Handle("/update/gauge/", http.StripPrefix("/update/gauge/", handlers.MetricHandler{
		MetricType: metrics.Gauge,
		Storage:    storage,
	}))

	return http.ListenAndServe("localhost:8080", mux)
}

func main() {
	if err := runServer(); err != nil {
		panic(err)
	}
}
