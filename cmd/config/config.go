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
	ConfigFilePath string `env:"CONFIG"`
	Delays         []time.Duration
	privateKey     *rsa.PrivateKey
	ServerAddr     NetAddress `env:"ADDRESS" json:"address"`
	StoreInterval  Duration   `env:"STORE_INTERVAL" json:"store_interval"`
	NeededRestore  bool       `env:"RESTORE" json:"restore"`
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

var defaultAgentConfig = AgentConfig{
	ServerAddr: NetAddress{
		Host: "localhost",
		Port: 8080,
	},
	ReportInterval: Duration{10 * time.Second},
	PollInterval:   Duration{2 * time.Second},
	RateLimit:      1,
	Key:            "",
	Delays:         []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
	CertFile:       "",
}

// AgentConfig структура для конфигурации агента сбора метрик.
type AgentConfig struct {
	CertFile       string     `env:"CRYPTO_KEY"`
	Key            string     `env:"KEY"`
	ConfigFilePath string     `env:"CONFIG"`
	ServerAddr     NetAddress `env:"ADDRESS"`
	Delays         []time.Duration
	ReportInterval Duration `env:"REPORT_INTERVAL"`
	PollInterval   Duration `env:"POLL_INTERVAL"`
	RateLimit      int      `env:"RATE_LIMIT"`
}

func (c *AgentConfig) parseFlags(programName string, args []string) error {
	flagSet := flag.NewFlagSet(programName, flag.ContinueOnError)

	flagSet.Var(&c.ServerAddr, "a", "address and port server")
	flagSet.DurationVar(&c.ReportInterval.Duration, "r", c.ReportInterval.Duration, "interval in seconds of sending metrics to server")
	flagSet.DurationVar(&c.PollInterval.Duration, "p", c.PollInterval.Duration, "interval in seconds of polling metrics")
	flagSet.StringVar(&c.Key, "k", c.Key, "key for signing queries")
	flagSet.IntVar(&c.RateLimit, "l", c.RateLimit, "rate limit outgoing requests to the server")
	flagSet.StringVar(&c.CertFile, "crypto-key", c.CertFile, "path to cert file")

	err := flagSet.Parse(args)
	if err != nil {
		return fmt.Errorf("parse config from flags: %w", err)
	}
	return nil
}

func (c *AgentConfig) parseEnv() error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("parse config from env: %w", err)
	}

	return nil
}

func (c *AgentConfig) parseFile(path string) error {
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

var agentConfig *AgentConfig = nil

// GetAgentConfig возвращает конфигурацию агента на основании указанных флагов при запуске или указанных переменных окружения.
func GetAgentConfig() (*AgentConfig, error) {
	var err error
	sync.OnceFunc(func() {
		agentConfig, err = getAgentConfig(os.Args[0], os.Args[1:])
	})()
	return agentConfig, err
}

func getAgentConfig(programName string, args []string) (*AgentConfig, error) {
	tmpCfg := defaultAgentConfig

	if err := tmpCfg.parseFlags(programName, args); err != nil {
		return nil, err
	}

	if err := tmpCfg.parseEnv(); err != nil {
		return nil, err
	}

	cfg := defaultAgentConfig
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
