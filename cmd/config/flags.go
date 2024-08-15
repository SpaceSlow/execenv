package config

import (
	"flag"
	"time"
)

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
