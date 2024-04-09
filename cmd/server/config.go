package main

import (
	"github.com/SpaceSlow/execenv/cmd/flags"
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr flags.NetAddress `env:"ADDRESS"`
}

func GetConfigWithFlags() (*Config, error) {
	parseFlags()
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if cfg.ServerAddr.String() == "" {
		cfg.ServerAddr = flagRunAddr
	}

	return cfg, nil
}
