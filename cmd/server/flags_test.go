package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SpaceSlow/execenv/cmd/flags"
)

type wantFlags struct {
	RunAddr       flags.NetAddress
	StoreInterval uint
	StoragePath   string
	NeedRestore   bool
	DatabaseDSN   string
	Key           string
}

var standardFlags = wantFlags{
	RunAddr: flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	},
	StoreInterval: 300,
	StoragePath:   "/tmp/metrics-db.json",
	NeedRestore:   true,
	DatabaseDSN:   "",
	Key:           "",
}

func Test_parseFlags(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantFlags wantFlags
	}{
		{
			name: "checking standard flag values",
			args: nil,
			wantFlags: wantFlags{
				RunAddr:       standardFlags.RunAddr,
				StoreInterval: standardFlags.StoreInterval,
				StoragePath:   standardFlags.StoragePath,
				NeedRestore:   standardFlags.NeedRestore,
				DatabaseDSN:   standardFlags.DatabaseDSN,
				Key:           standardFlags.Key,
			},
		},
		{
			name: "checking non-standard port flag",
			args: []string{
				"-a=:8081",
			},
			wantFlags: wantFlags{
				RunAddr: flags.NetAddress{
					Host: "",
					Port: 8081,
				},
				StoreInterval: standardFlags.StoreInterval,
				StoragePath:   standardFlags.StoragePath,
				NeedRestore:   standardFlags.NeedRestore,
				DatabaseDSN:   standardFlags.DatabaseDSN,
				Key:           standardFlags.Key,
			},
		},
		{
			name: "checking sync adding into file flag",
			args: []string{
				"-i=0",
			},
			wantFlags: wantFlags{
				RunAddr:       standardFlags.RunAddr,
				StoreInterval: 0,
				StoragePath:   standardFlags.StoragePath,
				NeedRestore:   standardFlags.NeedRestore,
				DatabaseDSN:   standardFlags.DatabaseDSN,
				Key:           standardFlags.Key,
			},
		},
		{
			name: "checking custom file storage flag",
			args: []string{
				"-f=/tmp/file",
			},
			wantFlags: wantFlags{
				RunAddr:       standardFlags.RunAddr,
				StoreInterval: standardFlags.StoreInterval,
				StoragePath:   "/tmp/file",
				NeedRestore:   standardFlags.NeedRestore,
				DatabaseDSN:   standardFlags.DatabaseDSN,
				Key:           standardFlags.Key,
			},
		},
		{
			name: "checking restore flag",
			args: []string{
				"-r",
			},
			wantFlags: wantFlags{
				RunAddr:       standardFlags.RunAddr,
				StoreInterval: standardFlags.StoreInterval,
				StoragePath:   standardFlags.StoragePath,
				NeedRestore:   true,
				DatabaseDSN:   standardFlags.DatabaseDSN,
				Key:           standardFlags.Key,
			},
		},
		{
			name: "checking database dsn flag",
			args: []string{
				"-d=postgres://username:password@localhost:5432/database_name",
			},
			wantFlags: wantFlags{
				RunAddr:       standardFlags.RunAddr,
				StoreInterval: standardFlags.StoreInterval,
				StoragePath:   standardFlags.StoragePath,
				NeedRestore:   standardFlags.NeedRestore,
				DatabaseDSN:   "postgres://username:password@localhost:5432/database_name",
				Key:           standardFlags.Key,
			},
		},
		{
			name: "checking setting non-empty key flag",
			args: []string{
				"-k=non-standard-key",
			},
			wantFlags: wantFlags{
				RunAddr:       standardFlags.RunAddr,
				StoreInterval: standardFlags.StoreInterval,
				StoragePath:   standardFlags.StoragePath,
				NeedRestore:   standardFlags.NeedRestore,
				DatabaseDSN:   standardFlags.DatabaseDSN,
				Key:           "non-standard-key",
			},
		},
		{
			name: "checking all flags",
			args: []string{
				"-a=example.com:80",
				"-i=10",
				"-r",
				"-f=/tmp/some-file.json",
				"-k=non-standard-key",
				"-d=postgres://username:password@localhost:5432/database_name",
			},
			wantFlags: wantFlags{
				RunAddr: flags.NetAddress{
					Host: "example.com",
					Port: 80,
				},
				StoreInterval: 10,
				StoragePath:   "/tmp/some-file.json",
				NeedRestore:   true,
				DatabaseDSN:   "postgres://username:password@localhost:5432/database_name",
				Key:           "non-standard-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseFlags("program", tt.args)

			if !assert.ObjectsAreEqual(tt.wantFlags.RunAddr, flagRunAddr) {
				t.Errorf("expected flagRunAddr: %v, got: %v", tt.wantFlags.RunAddr, flagRunAddr)
			}
			assert.Equalf(t, tt.wantFlags.StoragePath, flagStoragePath, `expected flagStoragePath: %v, got: %v`, tt.wantFlags.StoragePath, flagStoragePath)
			assert.Equalf(t, tt.wantFlags.StoreInterval, flagStoreInterval, `expected flagStoreInterval: %v, got: %v`, tt.wantFlags.StoreInterval, flagStoreInterval)
			assert.Equalf(t, tt.wantFlags.NeedRestore, flagNeedRestore, `expected flagNeedRestore: %v, got: %v`, tt.wantFlags.NeedRestore, flagNeedRestore)
			assert.Equalf(t, tt.wantFlags.DatabaseDSN, flagDatabaseDSN, `expected flagDatabaseDSN: "%v", got: "%v"`, tt.wantFlags.DatabaseDSN, flagDatabaseDSN)
			assert.Equalf(t, tt.wantFlags.Key, flagKey, `expected flagKey: "%v", got: "%v"`, tt.wantFlags.Key, flagKey)
		})
	}
}
