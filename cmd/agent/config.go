package main

import (
	"time"

	"github.com/caarlos0/env"

	"github.com/SpaceSlow/execenv/cmd/config"
)

type Config struct {
	Key            string            `env:"KEY"`
	ServerAddr     config.NetAddress `env:"ADDRESS"`
	Delays         []time.Duration
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
	RateLimit      int `env:"RATE_LIMIT"`
}

func GetConfigWithFlags(programName string, args []string) (*Config, error) {
	parseFlags(programName, args)
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
