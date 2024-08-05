package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerConfigWithFlags(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		want    *ServerConfig
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
				"STORE_INTERVAL":    "100",
				"RESTORE":           "false",
				"DATABASE_DSN":      "postgres://env:env@localhost:5432/env",
				"KEY":               "env",
			},
			flags: []string{"-a=:8080", "-f=/tmp/flag", "-i=0", "-r", "-d=postgres://flag:flag@localhost:5432/flag", "-k=flag"},
			want: &ServerConfig{
				StoragePath:   "/tmp/env",
				DatabaseDSN:   "postgres://env:env@localhost:5432/env",
				Key:           "env",
				Delays:        defaultServerConfig.Delays,
				ServerAddr:    NetAddress{Host: "", Port: 9090},
				StoreInterval: 100,
				NeededRestore: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			got, err := getServerConfigWithFlags("program", tt.flags)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
