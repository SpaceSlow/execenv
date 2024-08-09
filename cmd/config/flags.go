package config

import "flag"

var (
	flagServerRunAddr       NetAddress
	flagServerStoreInterval uint
	flagServerStoragePath   string
	flagServerNeedRestore   bool
	flagServerDatabaseDSN   string
	flagServerKey           string
	flagServerCertFile      string
)

func parseServerFlags(programName string, args []string) {
	flagSet := flag.NewFlagSet(programName, flag.ContinueOnError)

	flagServerRunAddr = defaultServerConfig.ServerAddr
	flagSet.Var(&flagServerRunAddr, "a", "address and port to run server")
	flagSet.UintVar(&flagServerStoreInterval, "i", defaultServerConfig.StoreInterval, "store interval in secs (default 300 sec)")
	flagSet.StringVar(&flagServerStoragePath, "f", defaultServerConfig.StoragePath, "file storage path (default /tmp/metrics-db.json")
	flagSet.BoolVar(&flagServerNeedRestore, "r", defaultServerConfig.NeededRestore, "needed loading saved metrics from file (default true)")
	flagSet.StringVar(&flagServerDatabaseDSN, "d", defaultServerConfig.DatabaseDSN, "PostgreSQL (ver. >=10) database DSN (example: postgres://username:password@localhost:5432/database_name")
	flagSet.StringVar(&flagServerKey, "k", defaultServerConfig.Key, "key for signing queries")
	flagSet.StringVar(&flagServerCertFile, "crypto-key", defaultServerConfig.CertFile, "path to cert file")

	flagSet.Parse(args)
}

var (
	flagAgentServerAddr     NetAddress
	flagAgentReportInterval int
	flagAgentPollInterval   int
	flagAgentKey            string
	flagAgentRateLimit      int
	flagAgentCertFile       string
)

func parseAgentFlags(programName string, args []string) {
	flagSet := flag.NewFlagSet(programName, flag.ContinueOnError)

	flagAgentServerAddr = defaultAgentConfig.ServerAddr
	flagSet.Var(&flagAgentServerAddr, "a", "address and port server")
	flagSet.IntVar(&flagAgentReportInterval, "r", defaultAgentConfig.ReportInterval, "interval in seconds of sending metrics to server")
	flagSet.IntVar(&flagAgentPollInterval, "p", defaultAgentConfig.PollInterval, "interval in seconds of polling metrics")
	flagSet.StringVar(&flagAgentKey, "k", defaultAgentConfig.Key, "key for signing queries")
	flagSet.IntVar(&flagAgentRateLimit, "l", defaultAgentConfig.RateLimit, "rate limit outgoing requests to the server")
	flagSet.StringVar(&flagAgentCertFile, "crypto-key", defaultAgentConfig.CertFile, "path to cert file")

	flagSet.Parse(args)
}
