package config

import "flag"

var (
	flagRunAddr       NetAddress
	flagStoreInterval uint
	flagStoragePath   string
	flagNeedRestore   bool
	flagDatabaseDSN   string
	flagKey           string
)

func parseFlags(programName string, args []string) {
	flagSet := flag.NewFlagSet(programName, flag.ContinueOnError)

	flagRunAddr = defaultConfig.ServerAddr
	flagSet.Var(&flagRunAddr, "a", "address and port to run server")
	flagSet.UintVar(&flagStoreInterval, "i", defaultConfig.StoreInterval, "store interval in secs (default 300 sec)")
	flagSet.StringVar(&flagStoragePath, "f", defaultConfig.StoragePath, "file storage path (default /tmp/metrics-db.json")
	flagSet.BoolVar(&flagNeedRestore, "r", defaultConfig.NeededRestore, "needed loading saved metrics from file (default true)")
	flagSet.StringVar(&flagDatabaseDSN, "d", defaultConfig.DatabaseDSN, "PostgreSQL (ver. >=10) database DSN (example: postgres://username:password@localhost:5432/database_name")
	flagSet.StringVar(&flagKey, "k", defaultConfig.Key, "key for signing queries")

	flagSet.Parse(args)
}
