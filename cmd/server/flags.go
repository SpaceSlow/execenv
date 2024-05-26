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
	flagDatabaseDSN   string
	flagKey           string
	flagRateLimit     int
)

func parseFlags() {
	flagRunAddr = flags.NetAddress{Host: "localhost", Port: 8080}
	flag.Var(&flagRunAddr, "a", "address and port to run server")
	flag.UintVar(&flagStoreInterval, "i", 300, "store interval in secs (default 300 sec)")
	flag.StringVar(&flagStoragePath, "f", "/tmp/metrics-db.json", "file storage path (default /tmp/metrics-db.json")
	flag.BoolVar(&flagNeedRestore, "r", true, "needed loading saved metrics from file (default true)")
	flag.StringVar(&flagDatabaseDSN, "d", "", "PostgreSQL (ver. >=10) database DSN (example: postgres://username:password@localhost:5432/database_name")
	flag.StringVar(&flagKey, "k", "", "key for signing queries")
	flag.IntVar(&flagRateLimit, "l", 1, "rate limit incoming requests to the server")

	flag.Parse()
}
