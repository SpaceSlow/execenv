package main

import (
	"flag"

	"github.com/SpaceSlow/execenv/cmd/flags"
)

var (
	flagServerAddr     flags.NetAddress
	flagReportInterval int
	flagPollInterval   int
	flagKey            string
	flagRateLimit      int
)

func parseFlags() {
	flagServerAddr = flags.NetAddress{Host: "localhost", Port: 8080}
	flag.Var(&flagServerAddr, "a", "address and port server")
	flag.IntVar(&flagReportInterval, "r", 10, "interval in seconds of sending metrics to server")
	flag.IntVar(&flagPollInterval, "p", 2, "interval in seconds of polling metrics")
	flag.StringVar(&flagKey, "k", "", "key for signing queries")
	flag.IntVar(&flagRateLimit, "l", 1, "rate limit outgoing requests to the server")

	flag.Parse()
}
