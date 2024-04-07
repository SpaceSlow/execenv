package main

import (
	"flag"

	"github.com/SpaceSlow/execenv/cmd/flags"
)

var (
	flagRunAddr flags.NetAddress
)

func parseFlags() {
	flagRunAddr = flags.NetAddress{Host: "localhost", Port: 8080}
	flag.Var(&flagRunAddr, "a", "address and port to run server")

	flag.Parse()
}
