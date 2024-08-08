package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/caarlos0/env"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// PrintBuildInfo выводит информацию о сборке.
// Необходима сборка с флагом ldflags следующих переменных:
//
// - github.com/SpaceSlow/execenv/cmd/config.buildVersion
//
// - github.com/SpaceSlow/execenv/cmd/config.buildDate
//
// - github.com/SpaceSlow/execenv/cmd/config.buildCommit
func PrintBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
}

var defaultServerConfig = &ServerConfig{
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

// ServerConfig структура для конфигурации сервера сбора метрик.
type ServerConfig struct {
	StoragePath   string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
	Key           string `env:"KEY"`
	Delays        []time.Duration
	ServerAddr    NetAddress `env:"ADDRESS"`
	StoreInterval uint       `env:"STORE_INTERVAL"`
	NeededRestore bool       `env:"RESTORE"`
}

var serverConfig *ServerConfig = nil

// GetServerConfig возвращает конфигурацию сервера на основании указанных флагов при запуске или указанных переменных окружения.
func GetServerConfig() (*ServerConfig, error) {
	var err error
	sync.OnceFunc(func() {
		serverConfig, err = getServerConfigWithFlags(os.Args[0], os.Args[1:])
	})()
	return serverConfig, err
}

func getServerConfigWithFlags(programName string, args []string) (*ServerConfig, error) {
	parseServerFlags(programName, args)
	cfg := &ServerConfig{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config from env: %w", err)
	}

	if cfg.ServerAddr.String() == "" {
		cfg.ServerAddr = flagServerRunAddr
	}
	if _, ok := os.LookupEnv("STORE_INTERVAL"); !ok {
		cfg.StoreInterval = flagServerStoreInterval
	}
	if _, ok := os.LookupEnv("FILE_STORAGE_PATH"); !ok {
		cfg.StoragePath = flagServerStoragePath
	}
	if _, ok := os.LookupEnv("RESTORE"); !ok {
		cfg.NeededRestore = flagServerNeedRestore
	}

	if _, ok := os.LookupEnv("DATABASE_DSN"); !ok {
		cfg.DatabaseDSN = flagServerDatabaseDSN
	}

	if cfg.Key == "" {
		cfg.Key = flagServerKey
	}

	cfg.Delays = defaultServerConfig.Delays

	return cfg, nil
}

var defaultAgentConfig = &AgentConfig{
	ServerAddr: NetAddress{
		Host: "localhost",
		Port: 8080,
	},
	ReportInterval: 10,
	PollInterval:   2,
	RateLimit:      1,
	Key:            "",
	Delays:         []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
}

// AgentConfig структура для конфигурации агента сбора метрик.
type AgentConfig struct {
	Key            string     `env:"KEY"`
	ServerAddr     NetAddress `env:"ADDRESS"`
	Delays         []time.Duration
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
	RateLimit      int `env:"RATE_LIMIT"`
}

var agentConfig *AgentConfig = nil

// GetAgentConfig возвращает конфигурацию агента на основании указанных флагов при запуске или указанных переменных окружения.
func GetAgentConfig() (*AgentConfig, error) {
	var err error
	sync.OnceFunc(func() {
		agentConfig, err = getAgentConfigWithFlags(os.Args[0], os.Args[1:])
	})()
	return agentConfig, err
}

func getAgentConfigWithFlags(programName string, args []string) (*AgentConfig, error) {
	parseAgentFlags(programName, args)
	cfg := &AgentConfig{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config from env: %w", err)
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = defaultAgentConfig.ReportInterval
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = defaultAgentConfig.PollInterval
	}
	if cfg.ServerAddr.String() == "" {
		cfg.ServerAddr = defaultAgentConfig.ServerAddr
	}
	if cfg.Key == "" {
		cfg.Key = defaultAgentConfig.Key
	}

	cfg.Delays = []time.Duration{
		time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	if cfg.RateLimit == 0 {
		cfg.RateLimit = defaultAgentConfig.RateLimit
	}

	return cfg, nil
}
