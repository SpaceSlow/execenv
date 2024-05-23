package main

import (
	"github.com/SpaceSlow/execenv/cmd/flags"
	"github.com/caarlos0/env"
	"time"
)

type Config struct {
	ServerAddr     flags.NetAddress `env:"ADDRESS"`
	ReportInterval int              `env:"REPORT_INTERVAL"`
	PollInterval   int              `env:"POLL_INTERVAL"`
	Key            string           `env:"KEY"`
	RateLimit      int              `env:"RATE_LIMIT"`
	Delays         []time.Duration
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
	if cfg.Key == "" {
		cfg.Key = flagKey
	}

	cfg.Delays = []time.Duration{
		time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	if cfg.RateLimit == 0 {
		cfg.RateLimit = flagRateLimit
	}

	return cfg, nil
}
