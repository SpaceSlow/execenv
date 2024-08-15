package config

import (
	"flag"
	"fmt"
	"time"
)

func ParseFlagsServerConfig(programName string, args []string, defaultCfg *ServerConfig) (*ServerConfig, error) {
	flagSet := flag.NewFlagSet(programName, flag.ContinueOnError)

	var cfg ServerConfig
	if defaultCfg != nil {
		cfg = *defaultCfg
	} else {
		cfg = defaultServerConfig
	}

	flagSet.Var(&cfg.ServerAddr, "a", "address and port to run server")
	flagSet.DurationVar(&cfg.StoreInterval.Duration, "i", cfg.StoreInterval.Duration, "store interval in secs (default 300 sec)")
	flagSet.StringVar(&cfg.StoragePath, "f", cfg.StoragePath, "file storage path (default /tmp/metrics-db.json")
	flagSet.BoolVar(&cfg.NeededRestore, "r", cfg.NeededRestore, "needed loading saved metrics from file (default true)")
	flagSet.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "PostgreSQL (ver. >=10) database DSN (example: postgres://username:password@localhost:5432/database_name")
	flagSet.StringVar(&cfg.Key, "k", cfg.Key, "key for signing queries")
	flagSet.StringVar(&cfg.CertFile, "crypto-key", cfg.CertFile, "path to cert file")

	flagSet.StringVar(&cfg.ConfigFilePath, "c", cfg.ConfigFilePath, "config file path")
	flagSet.StringVar(&cfg.ConfigFilePath, "config", cfg.ConfigFilePath, "config file path")

	err := flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse config from flags: %w", err)
	}
	return &cfg, nil
}

var (
	flagAgentServerAddr     NetAddress
	flagAgentReportInterval time.Duration
	flagAgentPollInterval   time.Duration
	flagAgentKey            string
	flagAgentRateLimit      int
	flagAgentCertFile       string
)

func parseAgentFlags(programName string, args []string) {
	flagSet := flag.NewFlagSet(programName, flag.ContinueOnError)

	flagAgentServerAddr = defaultAgentConfig.ServerAddr
	flagSet.Var(&flagAgentServerAddr, "a", "address and port server")
	flagSet.DurationVar(&flagAgentReportInterval, "r", defaultAgentConfig.ReportInterval, "interval in seconds of sending metrics to server")
	flagSet.DurationVar(&flagAgentPollInterval, "p", defaultAgentConfig.PollInterval, "interval in seconds of polling metrics")
	flagSet.StringVar(&flagAgentKey, "k", defaultAgentConfig.Key, "key for signing queries")
	flagSet.IntVar(&flagAgentRateLimit, "l", defaultAgentConfig.RateLimit, "rate limit outgoing requests to the server")
	flagSet.StringVar(&flagAgentCertFile, "crypto-key", defaultAgentConfig.CertFile, "path to cert file")

	flagSet.Parse(args)
}
