package main

import (
	"log"
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/config"
	"github.com/SpaceSlow/execenv/cmd/middlewares"
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
	if err := RunServer(middlewareHandlers...); err != nil {
		log.Fatalf("Error occured: %s.\r\nExiting...", err)
	}
}
