package main

import (
	"github.com/SpaceSlow/execenv/cmd/flags"
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr     flags.NetAddress `env:"ADDRESS"`
	ReportInterval int              `env:"REPORT_INTERVAL"`
	PollInterval   int              `env:"POLL_INTERVAL"`
}

func GetConfigWithFlags() (*Config, error) {
	parseFlags()
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = flagReportInterval
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = flagPollInterval
	}
	if cfg.ServerAddr.String() == "" {
		cfg.ServerAddr = flagServerAddr
	}

	return cfg, nil
}
