package main

import (
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/server/handlers"
	"github.com/SpaceSlow/execenv/cmd/server/storages"
)

func runServer() error {
	storage := storages.NewMemStorage()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.DefaultHandler)
	mux.Handle("/update/", http.StripPrefix("/update/", handlers.MetricHandler{Storage: storage}))

	return http.ListenAndServe("localhost:8080", mux)
}

func main() {
	if err := runServer(); err != nil {
		panic(err)
	}
}
