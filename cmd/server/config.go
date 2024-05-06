package main

import (
	"os"

	"github.com/SpaceSlow/execenv/cmd/flags"
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr    flags.NetAddress `env:"ADDRESS"`
	StoreInterval uint             `env:"STORE_INTERVAL"`
	StoragePath   string           `env:"FILE_STORAGE_PATH"`
	NeededRestore bool             `env:"RESTORE"`
	DatabaseDSN   string           `env:"DATABASE_DSN"`
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
	if _, ok := os.LookupEnv("STORE_INTERVAL"); !ok {
		cfg.StoreInterval = flagStoreInterval
	}
	if _, ok := os.LookupEnv("FILE_STORAGE_PATH"); !ok {
		cfg.StoragePath = flagStoragePath
	}
	if _, ok := os.LookupEnv("RESTORE"); !ok {
		cfg.NeededRestore = flagNeedRestore
	}

	if _, ok := os.LookupEnv("DATABASE_DSN"); !ok {
		cfg.DatabaseDSN = flagDatabaseDSN
	}

	return cfg, nil
}
