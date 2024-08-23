package config

import (
	"os"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerConfig_parseFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantCfg ServerConfig
	}{
		{
			name:    "standard flag values",
			args:    nil,
			wantCfg: defaultServerConfig,
		},
		{
			name: "non-standard port flag",
			args: []string{
				"-a=:8081",
			},
			wantCfg: ServerConfig{
				ServerAddr: NetAddress{
					Host: "",
					Port: 8081,
				},
				StoreInterval: defaultServerConfig.StoreInterval,
				StoragePath:   defaultServerConfig.StoragePath,
				NeededRestore: defaultServerConfig.NeededRestore,
				DatabaseDSN:   defaultServerConfig.DatabaseDSN,
				Key:           defaultServerConfig.Key,
			},
		},
		{
			name: "sync adding into file flag",
			args: []string{
				"-i=0s",
			},
			wantCfg: ServerConfig{
				ServerAddr:    defaultServerConfig.ServerAddr,
				StoreInterval: Duration{0},
				StoragePath:   defaultServerConfig.StoragePath,
				NeededRestore: defaultServerConfig.NeededRestore,
				DatabaseDSN:   defaultServerConfig.DatabaseDSN,
				Key:           defaultServerConfig.Key,
			},
		},
		{
			name: "custom file storage flag",
			args: []string{
				"-f=/tmp/file",
			},
			wantCfg: ServerConfig{
				ServerAddr:    defaultServerConfig.ServerAddr,
				StoreInterval: defaultServerConfig.StoreInterval,
				StoragePath:   "/tmp/file",
				NeededRestore: defaultServerConfig.NeededRestore,
				DatabaseDSN:   defaultServerConfig.DatabaseDSN,
				Key:           defaultServerConfig.Key,
			},
		},
		{
			name: "restore flag",
			args: []string{
				"-r",
			},
			wantCfg: ServerConfig{
				ServerAddr:    defaultServerConfig.ServerAddr,
				StoreInterval: defaultServerConfig.StoreInterval,
				StoragePath:   defaultServerConfig.StoragePath,
				NeededRestore: true,
				DatabaseDSN:   defaultServerConfig.DatabaseDSN,
				Key:           defaultServerConfig.Key,
			},
		},
		{
			name: "database dsn flag",
			args: []string{
				"-d=postgres://username:password@localhost:5432/database_name",
			},
			wantCfg: ServerConfig{
				ServerAddr:    defaultServerConfig.ServerAddr,
				StoreInterval: defaultServerConfig.StoreInterval,
				StoragePath:   defaultServerConfig.StoragePath,
				NeededRestore: defaultServerConfig.NeededRestore,
				DatabaseDSN:   "postgres://username:password@localhost:5432/database_name",
				Key:           defaultServerConfig.Key,
			},
		},
		{
			name: "setting non-empty key flag",
			args: []string{
				"-k=non-standard-key",
			},
			wantCfg: ServerConfig{
				ServerAddr:    defaultServerConfig.ServerAddr,
				StoreInterval: defaultServerConfig.StoreInterval,
				StoragePath:   defaultServerConfig.StoragePath,
				NeededRestore: defaultServerConfig.NeededRestore,
				DatabaseDSN:   defaultServerConfig.DatabaseDSN,
				Key:           "non-standard-key",
			},
		},
		{
			name: "all flags",
			args: []string{
				"-a=example.com:80",
				"-i=10s",
				"-r",
				"-f=/tmp/some-file.json",
				"-k=non-standard-key",
				"-d=postgres://username:password@localhost:5432/database_name",
			},
			wantCfg: ServerConfig{
				ServerAddr: NetAddress{
					Host: "example.com",
					Port: 80,
				},
				StoreInterval: Duration{10 * time.Second},
				StoragePath:   "/tmp/some-file.json",
				NeededRestore: true,
				DatabaseDSN:   "postgres://username:password@localhost:5432/database_name",
				Key:           "non-standard-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := defaultServerConfig
			err := config.parseFlags("program", tt.args)
			require.NoError(t, err)

			if !assert.ObjectsAreEqual(tt.wantCfg.ServerAddr, config.ServerAddr) {
				t.Errorf("expected ServerAddr: %v, got: %v", tt.wantCfg.ServerAddr, config.ServerAddr)
			}
			assert.Equalf(t, tt.wantCfg.StoragePath, config.StoragePath, `expected StoragePath: %v, got: %v`, tt.wantCfg.StoragePath, config.StoragePath)
			assert.Equalf(t, tt.wantCfg.StoreInterval, config.StoreInterval, `expected StoreInterval: %v, got: %v`, tt.wantCfg.StoreInterval, config.StoreInterval)
			assert.Equalf(t, tt.wantCfg.NeededRestore, config.NeededRestore, `expected NeedRestore: %v, got: %v`, tt.wantCfg.NeededRestore, config.NeededRestore)
			assert.Equalf(t, tt.wantCfg.DatabaseDSN, config.DatabaseDSN, `expected DatabaseDSN: "%v", got: "%v"`, tt.wantCfg.DatabaseDSN, config.DatabaseDSN)
			assert.Equalf(t, tt.wantCfg.Key, config.Key, `expected Key: "%v", got: "%v"`, tt.wantCfg.Key, config.Key)
		})
	}
}

func Test_getServerConfig(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		flags   []string
		want    ServerConfig
		wantErr bool
	}{
		{
			name: "standard config",
			want: defaultServerConfig,
		},
		{
			name: "incorrect server address in envs",
			envs: map[string]string{
				"ADDRESS": ":-1",
			},
			wantErr: true,
		},
		{
			name: "env priority on flags",
			envs: map[string]string{
				"ADDRESS":           ":9090",
				"FILE_STORAGE_PATH": "/tmp/env",
				"STORE_INTERVAL":    "100s",
				"RESTORE":           "false",
				"DATABASE_DSN":      "postgres://env:env@localhost:5432/env",
				"KEY":               "env",
				"CRYPTO_KEY":        "/tmp/cert.env.pem",
			},
			flags: []string{"-a=:8080", "-f=/tmp/flag", "-i=0s", "-r", "-d=postgres://flag:flag@localhost:5432/flag", "-k=flag", "-crypto-key=/tmp/cert.flag.pem"},
			want: ServerConfig{
				StoragePath:    "/tmp/env",
				DatabaseDSN:    "postgres://env:env@localhost:5432/env",
				Key:            "env",
				Delays:         defaultServerConfig.Delays,
				ServerAddr:     NetAddress{Host: "", Port: 9090},
				StoreInterval:  Duration{100 * time.Second},
				NeededRestore:  false,
				PrivateKeyFile: "/tmp/cert.env.pem",
			},
		},
		{
			name:  "only flags",
			flags: []string{"-a=:8080", "-f=/tmp/flag", "-i=0s", "-r", "-d=postgres://flag:flag@localhost:5432/flag", "-k=flag", "-crypto-key=/tmp/cert.flag.pem"},
			want: ServerConfig{
				StoragePath:    "/tmp/flag",
				DatabaseDSN:    "postgres://flag:flag@localhost:5432/flag",
				Key:            "flag",
				Delays:         defaultServerConfig.Delays,
				ServerAddr:     NetAddress{Host: "", Port: 8080},
				StoreInterval:  Duration{0},
				NeededRestore:  true,
				PrivateKeyFile: "/tmp/cert.flag.pem",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			got, err := getServerConfig("program", tt.flags)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				assert.ObjectsAreEqual(tt.want, *got)
			}
		})
	}
}

func TestGetServerConfig(t *testing.T) {
	temp := slices.Clone(os.Args)
	os.Args = []string{"test"}

	firstCfg, err := GetServerConfig()
	require.NoError(t, err)
	secondCfg, err := GetServerConfig()
	require.NoError(t, err)

	assert.Same(t, secondCfg, firstCfg)

	os.Args = temp
}

func TestServerConfig_parseFile(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectedCfg ServerConfig
		expectedErr bool
	}{
		{
			name: "config with all filled field",
			data: []byte(`
				{
					"address": "localhost:8080", 
					"restore": true, 
					"store_interval": "1s",
					"store_file": "/path/to/file.db", 
					"database_dsn": "postgres://file:file@localhost:5432/file",
					"key": "key",
					"crypto_key": "/path/to/private.key"
				} 
			`),
			expectedCfg: ServerConfig{
				StoragePath:    "/path/to/file.db",
				Key:            "key",
				DatabaseDSN:    "postgres://file:file@localhost:5432/file",
				PrivateKeyFile: "/path/to/private.key",
				Delays:         nil,
				ServerAddr:     NetAddress{Host: "localhost", Port: 8080},
				StoreInterval:  Duration{time.Second},
				NeededRestore:  true,
			},
			expectedErr: false,
		},
		{
			name:        "incorrect config",
			data:        []byte(`incorrect config`),
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, err := generateConfigFile(tt.data)
			require.NoError(t, err)
			defer os.Remove(filename)

			cfg := ServerConfig{}
			err = cfg.parseFile(filename)
			assert.Equal(t, tt.expectedErr, err != nil)

			assert.Equal(t, tt.expectedCfg, cfg)
		})
	}
}

func Test_getAgentConfigWithFlags(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		want    *AgentConfig
		flags   []string
		wantErr bool
	}{
		{
			name: "standard config",
			want: &defaultAgentConfig,
		},
		{
			name: "incorrect server address in envs",
			envs: map[string]string{
				"ADDRESS": ":-1",
			},
			wantErr: true,
		},
		{
			name: "env priority on flags",
			envs: map[string]string{
				"ADDRESS":         ":9090",
				"REPORT_INTERVAL": "5s",
				"POLL_INTERVAL":   "1s",
				"RATE_LIMIT":      "10",
				"KEY":             "env",
				"CRYPTO_KEY":      "/tmp/cert.env.pem",
			},
			flags: []string{"-a=:8080", "-r=55s", "-p=11s", "-l=100", "-k=flag", "-crypto-key=/tmp/cert.flag.pem"},
			want: &AgentConfig{
				ServerAddr:     NetAddress{Host: "", Port: 9090},
				ReportInterval: Duration{5 * time.Second},
				PollInterval:   Duration{1 * time.Second},
				RateLimit:      10,
				Key:            "env",
				Delays:         defaultServerConfig.Delays,
				CertFile:       "/tmp/cert.env.pem",
			},
		},
		{
			name:  "only flags",
			flags: []string{"-a=:8080", "-r=55s", "-p=11s", "-l=100", "-k=flag", "-crypto-key=/tmp/cert.flag.pem"},
			want: &AgentConfig{
				ServerAddr:     NetAddress{Host: "", Port: 8080},
				ReportInterval: Duration{55 * time.Second},
				PollInterval:   Duration{11 * time.Second},
				RateLimit:      100,
				Key:            "flag",
				Delays:         defaultServerConfig.Delays,
				CertFile:       "/tmp/cert.flag.pem",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			got, err := getAgentConfig("program", tt.flags)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

var standardFlags = AgentConfig{
	ServerAddr: NetAddress{
		Host: "localhost",
		Port: 8080,
	},
	ReportInterval: Duration{10 * time.Second},
	PollInterval:   Duration{2 * time.Second},
	RateLimit:      1,
	Key:            "",
}

func TestAgentConfig_parseFlags(t *testing.T) {
	tests := []struct {
		args    []string
		name    string
		wantCfg AgentConfig
	}{
		{
			name: "checking standard flag values",
			args: nil,
			wantCfg: AgentConfig{
				ServerAddr:     standardFlags.ServerAddr,
				ReportInterval: standardFlags.ReportInterval,
				PollInterval:   standardFlags.PollInterval,
				RateLimit:      standardFlags.RateLimit,
				Key:            standardFlags.Key,
			},
		},
		{
			name: "checking non-standard port flag",
			args: []string{
				"-a=:8081",
			},
			wantCfg: AgentConfig{
				ServerAddr: NetAddress{
					Host: "",
					Port: 8081,
				},
				ReportInterval: standardFlags.ReportInterval,
				PollInterval:   standardFlags.PollInterval,
				RateLimit:      standardFlags.RateLimit,
				Key:            standardFlags.Key,
			},
		},
		{
			name: "checking non-standard reporting interval flag",
			args: []string{
				"-r=30s",
			},
			wantCfg: AgentConfig{
				ServerAddr:     standardFlags.ServerAddr,
				ReportInterval: Duration{30 * time.Second},
				PollInterval:   standardFlags.PollInterval,
				RateLimit:      standardFlags.RateLimit,
				Key:            standardFlags.Key,
			},
		},
		{
			name: "checking non-standard polling interval flag",
			args: []string{
				"-p=2s",
			},
			wantCfg: AgentConfig{
				ServerAddr:     standardFlags.ServerAddr,
				ReportInterval: standardFlags.ReportInterval,
				PollInterval:   Duration{2 * time.Second},
				RateLimit:      standardFlags.RateLimit,
				Key:            standardFlags.Key,
			},
		},
		{
			name: "checking non-standard rate limit flag",
			args: []string{
				"-l=3",
			},
			wantCfg: AgentConfig{
				ServerAddr:     standardFlags.ServerAddr,
				ReportInterval: standardFlags.ReportInterval,
				PollInterval:   standardFlags.PollInterval,
				RateLimit:      3,
				Key:            standardFlags.Key,
			},
		},
		{
			name: "checking setting non-empty key flag",
			args: []string{
				"-k=non-standard-key",
			},
			wantCfg: AgentConfig{
				ServerAddr:     standardFlags.ServerAddr,
				ReportInterval: standardFlags.ReportInterval,
				PollInterval:   standardFlags.PollInterval,
				RateLimit:      standardFlags.RateLimit,
				Key:            "non-standard-key",
			},
		},
		{
			name: "checking all flags",
			args: []string{
				"-a=example.com:80",
				"-l=10",
				"-r=5s",
				"-p=1s",
				"-k=non-standard-key",
			},
			wantCfg: AgentConfig{
				ServerAddr: NetAddress{
					Host: "example.com",
					Port: 80,
				},
				ReportInterval: Duration{5 * time.Second},
				PollInterval:   Duration{1 * time.Second},
				RateLimit:      10,
				Key:            "non-standard-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := defaultAgentConfig
			err := config.parseFlags("program", tt.args)
			require.NoError(t, err)

			if !assert.ObjectsAreEqual(tt.wantCfg.ServerAddr, config.ServerAddr) {
				t.Errorf("expected ServerAddr: %v, got: %v", tt.wantCfg.ServerAddr, config.ServerAddr)
			}
			assert.Equalf(t, tt.wantCfg.ReportInterval, config.ReportInterval, `expected ReportInterval: %v, got: %v`, tt.wantCfg.ReportInterval, config.ReportInterval)
			assert.Equalf(t, tt.wantCfg.PollInterval, config.PollInterval, `expected PollInterval: %v, got: %v`, tt.wantCfg.PollInterval, config.PollInterval)
			assert.Equalf(t, tt.wantCfg.RateLimit, config.RateLimit, `expected RateLimit: %v, got: %v`, tt.wantCfg.RateLimit, config.RateLimit)
			assert.Equalf(t, tt.wantCfg.Key, config.Key, `expected Key: "%v", got: "%v"`, tt.wantCfg.Key, config.Key)
		})
	}
}

func TestGetAgentConfig(t *testing.T) {
	temp := slices.Clone(os.Args)
	os.Args = []string{"test"}

	firstCfg, err := GetAgentConfig()
	require.NoError(t, err)
	secondCfg, err := GetAgentConfig()
	require.NoError(t, err)

	assert.Same(t, secondCfg, firstCfg)

	os.Args = temp
}

func TestAgentConfig_parseFile(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectedCfg AgentConfig
		expectedErr bool
	}{
		{
			name: "config with all filled field",
			data: []byte(`
				{
					"address": "localhost:8080",
					"report_interval": "1s",
					"poll_interval": "1s",
					"crypto_key": "/path/to/cert.pem",
					"rate_limit": 4,
					"key": "key"
				}  
			`),
			expectedCfg: AgentConfig{
				Key:            "key",
				CertFile:       "/path/to/cert.pem",
				ReportInterval: Duration{time.Second},
				PollInterval:   Duration{time.Second},
				RateLimit:      4,
				ServerAddr:     NetAddress{Host: "localhost", Port: 8080},
			},
			expectedErr: false,
		},
		{
			name:        "incorrect config",
			data:        []byte(`incorrect config`),
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, err := generateConfigFile(tt.data)
			require.NoError(t, err)
			defer os.Remove(filename)

			cfg := AgentConfig{}
			err = cfg.parseFile(filename)
			assert.Equal(t, tt.expectedErr, err != nil)

			assert.Equal(t, tt.expectedCfg, cfg)
		})
	}
}

func generateConfigFile(data []byte) (string, error) {
	file, err := os.CreateTemp(os.TempDir(), "temp.*.cfg")
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
