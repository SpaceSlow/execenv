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

func parseFlags(programName string, args []string) {
	flagSet := flag.NewFlagSet(programName, flag.ContinueOnError)

	flagServerAddr = flags.NetAddress{Host: "localhost", Port: 8080}
	flagSet.Var(&flagServerAddr, "a", "address and port server")
	flagSet.IntVar(&flagReportInterval, "r", 10, "interval in seconds of sending metrics to server")
	flagSet.IntVar(&flagPollInterval, "p", 2, "interval in seconds of polling metrics")
	flagSet.StringVar(&flagKey, "k", "", "key for signing queries")
	flagSet.IntVar(&flagRateLimit, "l", 1, "rate limit outgoing requests to the server")

	flagSet.Parse(args)
}
