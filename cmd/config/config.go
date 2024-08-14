package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
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
	StoreInterval: 300 * time.Second,
	StoragePath:   "/tmp/metrics-db.json",
	NeededRestore: true,
	DatabaseDSN:   "",
	Key:           "",
	CertFile:      "",
	Delays:        []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
}

// ServerConfig структура для конфигурации сервера сбора метрик.
type ServerConfig struct {
	StoragePath    string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics-db.json" json:"store_file"`
	DatabaseDSN    string `env:"DATABASE_DSN" envDefault:"" json:"database_dsn"`
	Key            string `env:"KEY" envDefault:""`
	CertFile       string `env:"CRYPTO_KEY" envDefault:"" json:"crypto_key"`
	Delays         []time.Duration
	privateKey     *rsa.PrivateKey
	ServerAddr     NetAddress    `env:"ADDRESS" envDefault:"localhost:8080" json:"address"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"300s" json:"store_interval"`
	NeededRestore  bool          `env:"RESTORE" envDefault:"true" json:"restore"`
	ConfigFilePath string        `env:"CONFIG" envDefault:""`
}

func (c *ServerConfig) setPrivateKey() error {
	if c.CertFile == "" {
		return nil
	}
	keyBytes, err := os.ReadFile(c.CertFile)
	if err != nil {
		return err
	}
	keyBlock, _ := pem.Decode(keyBytes)
	c.privateKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return err
	}
	return nil
}

func (c *ServerConfig) PrivateKey() *rsa.PrivateKey {
	return c.privateKey
}

func (c *ServerConfig) UpdateDefaultFields(cfg *ServerConfig) {
	// TODO checking c.field == defaultConfig.field then update from cfg
}

var serverConfig *ServerConfig = nil

// GetServerConfig возвращает конфигурацию сервера на основании указанных флагов при запуске или указанных переменных окружения.
func GetServerConfig() (*ServerConfig, error) {
	var err error
	sync.OnceFunc(func() {
		serverConfig, err = getServerConfigWithFlags(os.Args[0], os.Args[1:])
		if err != nil {
			return
		}
		err = serverConfig.setPrivateKey()
	})()
	return serverConfig, err
}

func setServerFlagValues(cfg *ServerConfig) {
	if _, ok := os.LookupEnv("ADDRESS"); !ok {
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

	if _, ok := os.LookupEnv("KEY"); !ok {
		cfg.Key = flagServerKey
	}
	if _, ok := os.LookupEnv("CRYPTO_KEY"); !ok {
		cfg.CertFile = flagServerCertFile
	}

	if _, ok := os.LookupEnv("CONFIG"); !ok {
		cfg.ConfigFilePath = flagServerConfigFile
	}

	cfg.Delays = defaultServerConfig.Delays
}

func getServerConfigWithFlags(programName string, args []string) (*ServerConfig, error) {
	parseServerFlags(programName, args)
	cfg := &ServerConfig{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config from env: %w", err)
	}
	fCfg, err := ParseServerConfig(cfg.ConfigFilePath)
	if err != nil && !errors.Is(err, ErrEmptyPath) {
		return nil, err
	}
	setServerFlagValues(cfg)
	cfg.UpdateDefaultFields(fCfg)

	return cfg, nil
}

func ParseServerConfig(path string) (*ServerConfig, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := defaultServerConfig
	err = json.Unmarshal(data, cfg) // TODO fix time.Duration unmarshalling
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

var defaultAgentConfig = &AgentConfig{
	ServerAddr: NetAddress{
		Host: "localhost",
		Port: 8080,
	},
	ReportInterval: 10 * time.Second,
	PollInterval:   2 * time.Second,
	RateLimit:      1,
	Key:            "",
	Delays:         []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
	CertFile:       "",
}

// AgentConfig структура для конфигурации агента сбора метрик.
type AgentConfig struct {
	CertFile       string     `env:"CRYPTO_KEY"`
	Key            string     `env:"KEY"`
	ServerAddr     NetAddress `env:"ADDRESS"`
	Delays         []time.Duration
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	RateLimit      int           `env:"RATE_LIMIT"`
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

func setAgentDefaultValues(cfg *AgentConfig) {
	if _, ok := os.LookupEnv("REPORT_INTERVAL"); !ok {
		cfg.ReportInterval = flagAgentReportInterval
	}
	if _, ok := os.LookupEnv("POLL_INTERVAL"); !ok {
		cfg.PollInterval = flagAgentPollInterval
	}
	if _, ok := os.LookupEnv("ADDRESS"); !ok {
		cfg.ServerAddr = flagAgentServerAddr
	}
	if cfg.Key == "" {
		cfg.Key = flagAgentKey
	}

	cfg.Delays = defaultAgentConfig.Delays

	if cfg.RateLimit == 0 {
		cfg.RateLimit = flagAgentRateLimit
	}
	if cfg.CertFile == "" {
		cfg.CertFile = flagAgentCertFile
	}
}

func getAgentConfigWithFlags(programName string, args []string) (*AgentConfig, error) {
	parseAgentFlags(programName, args)
	cfg := &AgentConfig{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config from env: %w", err)
	}
	setAgentDefaultValues(cfg)
	return cfg, nil
}
