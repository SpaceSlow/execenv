package main

import (
	"fmt"
	"github.com/SpaceSlow/execenv/cmd/flags"
	"github.com/caarlos0/env"
	"os"
)

type Config struct {
	ServerAddr    flags.NetAddress `env:"ADDRESS"`
	StoreInterval int              `env:"STORE_INTERVAL"`
	StoragePath   string           `env:"FILE_STORAGE_PATH"`
	NeededRestore bool             `env:"RESTORE"`
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
	if os.Getenv("STORE_INTERVAL") == "" {
		cfg.StoreInterval = flagStoreInterval
	}
	if cfg.StoragePath == "" {
		cfg.StoragePath = flagStoragePath
	}
	if os.Getenv("RESTORE") == "" {
		cfg.NeededRestore = flagNeedRestore
	}

	fmt.Println(cfg)
	return cfg, nil
}
