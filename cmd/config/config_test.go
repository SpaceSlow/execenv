package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_getServerConfig(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		want    ServerConfig
		flags   []string
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
				StoragePath:   "/tmp/env",
				DatabaseDSN:   "postgres://env:env@localhost:5432/env",
				Key:           "env",
				Delays:        defaultServerConfig.Delays,
				ServerAddr:    NetAddress{Host: "", Port: 9090},
				StoreInterval: 100 * time.Second,
				NeededRestore: false,
				CertFile:      "/tmp/cert.env.pem",
			},
		},
		{
			name:  "only flags",
			flags: []string{"-a=:8080", "-f=/tmp/flag", "-i=0s", "-r", "-d=postgres://flag:flag@localhost:5432/flag", "-k=flag", "-crypto-key=/tmp/cert.flag.pem"},
			want: ServerConfig{
				StoragePath:   "/tmp/flag",
				DatabaseDSN:   "postgres://flag:flag@localhost:5432/flag",
				Key:           "flag",
				Delays:        defaultServerConfig.Delays,
				ServerAddr:    NetAddress{Host: "", Port: 8080},
				StoreInterval: 0,
				NeededRestore: true,
				CertFile:      "/tmp/cert.flag.pem",
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
			want: defaultAgentConfig,
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
				ReportInterval: 5 * time.Second,
				PollInterval:   1 * time.Second,
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
				ReportInterval: 55 * time.Second,
				PollInterval:   11 * time.Second,
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

			got, err := getAgentConfigWithFlags("program", tt.flags)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
