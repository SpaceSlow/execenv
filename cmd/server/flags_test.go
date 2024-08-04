package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SpaceSlow/execenv/cmd/flags"
)

type wantFlags struct {
	StoragePath   string
	DatabaseDSN   string
	Key           string
	ServerAddr    flags.NetAddress
	StoreInterval uint
	NeededRestore bool
}

func Test_parseFlags(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantFlags wantFlags
	}{
		{
			name: "standard flag values",
			args: nil,
			wantFlags: wantFlags{
				ServerAddr:    DefaultConfig.ServerAddr,
				StoreInterval: DefaultConfig.StoreInterval,
				StoragePath:   DefaultConfig.StoragePath,
				NeededRestore: DefaultConfig.NeededRestore,
				DatabaseDSN:   DefaultConfig.DatabaseDSN,
				Key:           DefaultConfig.Key,
			},
		},
		{
			name: "non-standard port flag",
			args: []string{
				"-a=:8081",
			},
			wantFlags: wantFlags{
				ServerAddr: flags.NetAddress{
					Host: "",
					Port: 8081,
				},
				StoreInterval: DefaultConfig.StoreInterval,
				StoragePath:   DefaultConfig.StoragePath,
				NeededRestore: DefaultConfig.NeededRestore,
				DatabaseDSN:   DefaultConfig.DatabaseDSN,
				Key:           DefaultConfig.Key,
			},
		},
		{
			name: "sync adding into file flag",
			args: []string{
				"-i=0",
			},
			wantFlags: wantFlags{
				ServerAddr:    DefaultConfig.ServerAddr,
				StoreInterval: 0,
				StoragePath:   DefaultConfig.StoragePath,
				NeededRestore: DefaultConfig.NeededRestore,
				DatabaseDSN:   DefaultConfig.DatabaseDSN,
				Key:           DefaultConfig.Key,
			},
		},
		{
			name: "custom file storage flag",
			args: []string{
				"-f=/tmp/file",
			},
			wantFlags: wantFlags{
				ServerAddr:    DefaultConfig.ServerAddr,
				StoreInterval: DefaultConfig.StoreInterval,
				StoragePath:   "/tmp/file",
				NeededRestore: DefaultConfig.NeededRestore,
				DatabaseDSN:   DefaultConfig.DatabaseDSN,
				Key:           DefaultConfig.Key,
			},
		},
		{
			name: "restore flag",
			args: []string{
				"-r",
			},
			wantFlags: wantFlags{
				ServerAddr:    DefaultConfig.ServerAddr,
				StoreInterval: DefaultConfig.StoreInterval,
				StoragePath:   DefaultConfig.StoragePath,
				NeededRestore: true,
				DatabaseDSN:   DefaultConfig.DatabaseDSN,
				Key:           DefaultConfig.Key,
			},
		},
		{
			name: "database dsn flag",
			args: []string{
				"-d=postgres://username:password@localhost:5432/database_name",
			},
			wantFlags: wantFlags{
				ServerAddr:    DefaultConfig.ServerAddr,
				StoreInterval: DefaultConfig.StoreInterval,
				StoragePath:   DefaultConfig.StoragePath,
				NeededRestore: DefaultConfig.NeededRestore,
				DatabaseDSN:   "postgres://username:password@localhost:5432/database_name",
				Key:           DefaultConfig.Key,
			},
		},
		{
			name: "setting non-empty key flag",
			args: []string{
				"-k=non-standard-key",
			},
			wantFlags: wantFlags{
				ServerAddr:    DefaultConfig.ServerAddr,
				StoreInterval: DefaultConfig.StoreInterval,
				StoragePath:   DefaultConfig.StoragePath,
				NeededRestore: DefaultConfig.NeededRestore,
				DatabaseDSN:   DefaultConfig.DatabaseDSN,
				Key:           "non-standard-key",
			},
		},
		{
			name: "all flags",
			args: []string{
				"-a=example.com:80",
				"-i=10",
				"-r",
				"-f=/tmp/some-file.json",
				"-k=non-standard-key",
				"-d=postgres://username:password@localhost:5432/database_name",
			},
			wantFlags: wantFlags{
				ServerAddr: flags.NetAddress{
					Host: "example.com",
					Port: 80,
				},
				StoreInterval: 10,
				StoragePath:   "/tmp/some-file.json",
				NeededRestore: true,
				DatabaseDSN:   "postgres://username:password@localhost:5432/database_name",
				Key:           "non-standard-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseFlags("program", tt.args)

			if !assert.ObjectsAreEqual(tt.wantFlags.ServerAddr, flagRunAddr) {
				t.Errorf("expected flagServerAddr: %v, got: %v", tt.wantFlags.ServerAddr, flagRunAddr)
			}
			assert.Equalf(t, tt.wantFlags.StoragePath, flagStoragePath, `expected flagStoragePath: %v, got: %v`, tt.wantFlags.StoragePath, flagStoragePath)
			assert.Equalf(t, tt.wantFlags.StoreInterval, flagStoreInterval, `expected flagStoreInterval: %v, got: %v`, tt.wantFlags.StoreInterval, flagStoreInterval)
			assert.Equalf(t, tt.wantFlags.NeededRestore, flagNeedRestore, `expected flagNeedRestore: %v, got: %v`, tt.wantFlags.NeededRestore, flagNeedRestore)
			assert.Equalf(t, tt.wantFlags.DatabaseDSN, flagDatabaseDSN, `expected flagDatabaseDSN: "%v", got: "%v"`, tt.wantFlags.DatabaseDSN, flagDatabaseDSN)
			assert.Equalf(t, tt.wantFlags.Key, flagKey, `expected flagKey: "%v", got: "%v"`, tt.wantFlags.Key, flagKey)
		})
	}
}
