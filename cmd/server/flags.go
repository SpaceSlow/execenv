package main

import (
	"flag"

	"github.com/SpaceSlow/execenv/cmd/flags"
)

var (
	flagRunAddr       flags.NetAddress
	flagStoreInterval uint
	flagStoragePath   string
	flagNeedRestore   bool
)

func parseFlags() {
	flagRunAddr = flags.NetAddress{Host: "localhost", Port: 8080}
	flag.Var(&flagRunAddr, "a", "address and port to run server")
	flag.UintVar(&flagStoreInterval, "i", 300, "store interval in secs (default 300 sec)")
	flag.StringVar(&flagStoragePath, "f", "/tmp/metrics-db.json", "file storage path (default /tmp/metrics-db.json")
	flag.BoolVar(&flagNeedRestore, "r", true, "needed loading saved metrics from file (default true)")

	flag.Parse()
}
