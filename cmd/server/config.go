package main

import (
	"os"
	"time"

	"github.com/caarlos0/env"

	"github.com/SpaceSlow/execenv/cmd/flags"
)

type Config struct {
	ServerAddr    flags.NetAddress `env:"ADDRESS"`
	StoreInterval uint             `env:"STORE_INTERVAL"`
	StoragePath   string           `env:"FILE_STORAGE_PATH"`
	NeededRestore bool             `env:"RESTORE"`
	DatabaseDSN   string           `env:"DATABASE_DSN"`
	Key           string           `env:"KEY"`
	Delays        []time.Duration
}

// GetConfigWithFlags возвращает конфигурацию сервера на основании указанных флагов при запуске или указанных переменных окружения.
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

	if cfg.Key == "" {
		cfg.Key = flagKey
	}

	cfg.Delays = []time.Duration{
		time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	return cfg, nil
}
