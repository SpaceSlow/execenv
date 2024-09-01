package main

import (
	"log"
	"net/http"

	"github.com/SpaceSlow/execenv/internal/config"
	"github.com/SpaceSlow/execenv/internal/middlewares"
	"github.com/SpaceSlow/execenv/internal/server"
)

func main() {
	config.PrintBuildInfo()
	middlewareHandlers := []func(next http.Handler) http.Handler{
		middlewares.WithSigning,
		middlewares.WithCompressing,
		middlewares.WithDecryption,
		middlewares.WithCheckingTrustedSubnet,
		middlewares.WithLogging,
	}
	if err := server.RunServer(middlewareHandlers...); err != nil {
		log.Fatalf("Error occured: %s.\r\nExiting...", err)
	}
}
