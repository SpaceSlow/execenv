package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type wantServerFlags struct {
	StoragePath   string
	DatabaseDSN   string
	Key           string
	ServerAddr    NetAddress
	StoreInterval Duration
	NeededRestore bool
}

func Test_ParseFlagsServerConfig(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantFlags wantServerFlags
	}{
		{
			name: "standard flag values",
			args: nil,
			wantFlags: wantServerFlags{
				ServerAddr:    defaultServerConfig.ServerAddr,
				StoreInterval: defaultServerConfig.StoreInterval,
				StoragePath:   defaultServerConfig.StoragePath,
				NeededRestore: defaultServerConfig.NeededRestore,
				DatabaseDSN:   defaultServerConfig.DatabaseDSN,
				Key:           defaultServerConfig.Key,
			},
		},
		{
			name: "non-standard port flag",
			args: []string{
				"-a=:8081",
			},
			wantFlags: wantServerFlags{
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
			wantFlags: wantServerFlags{
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
			wantFlags: wantServerFlags{
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
			wantFlags: wantServerFlags{
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
			wantFlags: wantServerFlags{
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
			wantFlags: wantServerFlags{
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
			wantFlags: wantServerFlags{
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
			config, err := ParseFlagsServerConfig("program", tt.args, nil)
			require.NoError(t, err)

			if !assert.ObjectsAreEqual(tt.wantFlags.ServerAddr, config.ServerAddr) {
				t.Errorf("expected flagServerAddr: %v, got: %v", tt.wantFlags.ServerAddr, config.ServerAddr)
			}
			assert.Equalf(t, tt.wantFlags.StoragePath, config.StoragePath, `expected flagStoragePath: %v, got: %v`, tt.wantFlags.StoragePath, config.StoragePath)
			assert.Equalf(t, tt.wantFlags.StoreInterval, config.StoreInterval, `expected flagStoreInterval: %v, got: %v`, tt.wantFlags.StoreInterval, config.StoreInterval)
			assert.Equalf(t, tt.wantFlags.NeededRestore, config.NeededRestore, `expected flagNeedRestore: %v, got: %v`, tt.wantFlags.NeededRestore, config.NeededRestore)
			assert.Equalf(t, tt.wantFlags.DatabaseDSN, config.DatabaseDSN, `expected flagDatabaseDSN: "%v", got: "%v"`, tt.wantFlags.DatabaseDSN, config.DatabaseDSN)
			assert.Equalf(t, tt.wantFlags.Key, config.Key, `expected flagKey: "%v", got: "%v"`, tt.wantFlags.Key, config.Key)
		})
	}
}

type wantAgentFlags struct {
	Key            string
	ServerAddr     NetAddress
	ReportInterval time.Duration
	PollInterval   time.Duration
	RateLimit      int
}

var standardFlags = wantAgentFlags{
	ServerAddr: NetAddress{
		Host: "localhost",
		Port: 8080,
	},
	ReportInterval: 10 * time.Second,
	PollInterval:   2 * time.Second,
	RateLimit:      1,
	Key:            "",
}

func Test_parseAgentFlags(t *testing.T) {
	tests := []struct {
		args      []string
		name      string
		wantFlags wantAgentFlags
	}{
		{
			name: "checking standard flag values",
			args: nil,
			wantFlags: wantAgentFlags{
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
			wantFlags: wantAgentFlags{
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
			wantFlags: wantAgentFlags{
				ServerAddr:     standardFlags.ServerAddr,
				ReportInterval: 30 * time.Second,
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
			wantFlags: wantAgentFlags{
				ServerAddr:     standardFlags.ServerAddr,
				ReportInterval: standardFlags.ReportInterval,
				PollInterval:   2 * time.Second,
				RateLimit:      standardFlags.RateLimit,
				Key:            standardFlags.Key,
			},
		},
		{
			name: "checking non-standard rate limit flag",
			args: []string{
				"-l=3",
			},
			wantFlags: wantAgentFlags{
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
			wantFlags: wantAgentFlags{
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
			wantFlags: wantAgentFlags{
				ServerAddr: NetAddress{
					Host: "example.com",
					Port: 80,
				},
				ReportInterval: 5 * time.Second,
				PollInterval:   1 * time.Second,
				RateLimit:      10,
				Key:            "non-standard-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseAgentFlags("program", tt.args)

			if !assert.ObjectsAreEqual(tt.wantFlags.ServerAddr, flagAgentServerAddr) {
				t.Errorf("expected flagServerAddr: %v, got: %v", tt.wantFlags.ServerAddr, flagAgentServerAddr)
			}
			assert.Equalf(t, tt.wantFlags.ReportInterval, flagAgentReportInterval, `expected flagReportInterval: %v, got: %v`, tt.wantFlags.ReportInterval, flagAgentReportInterval)
			assert.Equalf(t, tt.wantFlags.PollInterval, flagAgentPollInterval, `expected flagPollInterval: %v, got: %v`, tt.wantFlags.PollInterval, flagAgentPollInterval)
			assert.Equalf(t, tt.wantFlags.RateLimit, flagAgentRateLimit, `expected flagRateLimit: %v, got: %v`, tt.wantFlags.RateLimit, flagAgentRateLimit)
			assert.Equalf(t, tt.wantFlags.Key, flagAgentKey, `expected flagKey: "%v", got: "%v"`, tt.wantFlags.Key, flagAgentKey)
		})
	}
}
