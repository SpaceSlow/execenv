package main

import (
	"log"

	"github.com/SpaceSlow/execenv/internal/config"
	"github.com/SpaceSlow/execenv/internal/server"
)

func main() {
	config.PrintBuildInfo()
	if srv, err := server.NewServer(); err != nil || srv.Run() != nil {
		log.Fatalf("Error occured: %s.\r\nExiting...", err)
	}
}
