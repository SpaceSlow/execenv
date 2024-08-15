package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
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

var defaultServerConfig = ServerConfig{
	ServerAddr: NetAddress{
		Host: "localhost",
		Port: 8080,
	},
	StoreInterval: Duration{300 * time.Second},
	StoragePath:   "/tmp/metrics-db.json",
	NeededRestore: true,
	DatabaseDSN:   "",
	Key:           "",
	CertFile:      "",
	Delays:        []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
}

// ServerConfig структура для конфигурации сервера сбора метрик.
type ServerConfig struct {
	StoragePath    string `env:"FILE_STORAGE_PATH" json:"store_file"`
	DatabaseDSN    string `env:"DATABASE_DSN" json:"database_dsn"`
	Key            string `env:"KEY"`
	CertFile       string `env:"CRYPTO_KEY" json:"crypto_key"`
	Delays         []time.Duration
	privateKey     *rsa.PrivateKey
	ServerAddr     NetAddress `env:"ADDRESS" json:"address"`
	StoreInterval  Duration   `env:"STORE_INTERVAL" json:"store_interval"`
	NeededRestore  bool       `env:"RESTORE" json:"restore"`
	ConfigFilePath string     `env:"CONFIG"`
}

func (c *ServerConfig) parseFlags(programName string, args []string) error {
	flagSet := flag.NewFlagSet(programName, flag.ContinueOnError)

	flagSet.Var(&c.ServerAddr, "a", "address and port to run server")
	flagSet.DurationVar(&c.StoreInterval.Duration, "i", c.StoreInterval.Duration, "store interval in secs (default 300 sec)")
	flagSet.StringVar(&c.StoragePath, "f", c.StoragePath, "file storage path (default /tmp/metrics-db.json")
	flagSet.BoolVar(&c.NeededRestore, "r", c.NeededRestore, "needed loading saved metrics from file (default true)")
	flagSet.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "PostgreSQL (ver. >=10) database DSN (example: postgres://username:password@localhost:5432/database_name")
	flagSet.StringVar(&c.Key, "k", c.Key, "key for signing queries")
	flagSet.StringVar(&c.CertFile, "crypto-key", c.CertFile, "path to cert file")

	flagSet.StringVar(&c.ConfigFilePath, "c", c.ConfigFilePath, "config file path")
	flagSet.StringVar(&c.ConfigFilePath, "config", c.ConfigFilePath, "config file path")

	err := flagSet.Parse(args)
	if err != nil {
		return fmt.Errorf("parse config from flags: %w", err)
	}
	return nil
}

func (c *ServerConfig) parseEnv() error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("parse config from env: %w", err)
	}

	return nil
}

func (c *ServerConfig) parseFile(path string) error {
	if path == "" {
		return ErrEmptyPath
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("parse config from file: %w", err)
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("parse config from file: %w", err)
	}

	return nil
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

var serverConfig *ServerConfig = nil

// GetServerConfig возвращает конфигурацию сервера на основании указанных флагов при запуске или указанных переменных окружения.
func GetServerConfig() (*ServerConfig, error) {
	var err error
	sync.OnceFunc(func() {
		serverConfig, err = getServerConfig(os.Args[0], os.Args[1:])
		if err != nil {
			return
		}
		err = serverConfig.setPrivateKey()
	})()
	return serverConfig, err
}

func getServerConfig(programName string, args []string) (*ServerConfig, error) {
	tmpCfg := defaultServerConfig

	if err := tmpCfg.parseFlags(programName, args); err != nil {
		return nil, err
	}

	if err := tmpCfg.parseEnv(); err != nil {
		return nil, err
	}

	cfg := defaultServerConfig
	if err := cfg.parseFile(tmpCfg.ConfigFilePath); err != nil && !errors.Is(err, ErrEmptyPath) {
		return nil, err
	}

	if err := cfg.parseFlags(programName, args); err != nil {
		return nil, err
	}

	if err := cfg.parseEnv(); err != nil {
		return nil, err
	}

	return &cfg, nil
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
