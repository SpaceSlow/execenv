package main

import (
	"os"
	"time"

	"github.com/caarlos0/env"

	"github.com/SpaceSlow/execenv/cmd/flags"
)

var DefaultConfig = Config{
	ServerAddr: flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	},
	StoreInterval: 300,
	StoragePath:   "/tmp/metrics-db.json",
	NeededRestore: true,
	DatabaseDSN:   "",
	Key:           "",
	Delays:        []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
}

type Config struct {
	StoragePath   string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
	Key           string `env:"KEY"`
	Delays        []time.Duration
	ServerAddr    flags.NetAddress `env:"ADDRESS"`
	StoreInterval uint             `env:"STORE_INTERVAL"`
	NeededRestore bool             `env:"RESTORE"`
}

// GetConfigWithFlags возвращает конфигурацию сервера на основании указанных флагов при запуске или указанных переменных окружения.
func GetConfigWithFlags(programName string, args []string) (*Config, error) {
	parseFlags(programName, args)
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
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

	cfg.Delays = DefaultConfig.Delays

	return cfg, nil
}
