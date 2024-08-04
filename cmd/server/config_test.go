package main

import (
	"github.com/SpaceSlow/execenv/cmd/flags"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfigWithFlags(t *testing.T) {
	tests := []struct {
		err   error
		want  *Config
		name  string
		envs  map[string]string
		flags []string
	}{
		{
			name: "standard config",
			want: &DefaultConfig,
		},
		{
			name: "incorrect server address in envs",
			envs: map[string]string{
				"ADDRESS": ":-1",
			},
			err: flags.ErrIncorrectPort,
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
			want: &Config{
				StoragePath:   "/tmp/env",
				DatabaseDSN:   "postgres://env:env@localhost:5432/env",
				Key:           "env",
				Delays:        DefaultConfig.Delays,
				ServerAddr:    flags.NetAddress{Host: "", Port: 9090},
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

			got, err := GetConfigWithFlags("program", tt.flags)
			assert.Equal(t, err, tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
