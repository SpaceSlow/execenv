package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/caarlos0/env"
)

var defaultConfig = &Config{
	ServerAddr: NetAddress{
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
	ServerAddr    NetAddress `env:"ADDRESS"`
	StoreInterval uint       `env:"STORE_INTERVAL"`
	NeededRestore bool       `env:"RESTORE"`
}

var userConfig *Config = nil

// GetConfig возвращает конфигурацию сервера на основании указанных флагов при запуске или указанных переменных окружения.
func GetConfig() *Config {
	sync.OnceFunc(func() {
		var err error
		userConfig, err = getConfigWithFlags(os.Args[0], os.Args[1:])
		if err != nil {
			panic(err)
		}
	})()
	return userConfig
}

func getConfigWithFlags(programName string, args []string) (*Config, error) {
	parseFlags(programName, args)
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config from env: %w", err)
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

	cfg.Delays = defaultConfig.Delays

	return cfg, nil
}
