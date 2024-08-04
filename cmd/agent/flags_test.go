package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SpaceSlow/execenv/cmd/flags"
)

type wantFlags struct {
	ServerAddr     flags.NetAddress
	ReportInterval int
	PollInterval   int
	RateLimit      int
	Key            string
}

var standardFlags = wantFlags{
	ServerAddr: flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	},
	ReportInterval: 10,
	PollInterval:   2,
	RateLimit:      1,
	Key:            "",
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
			wantFlags: wantFlags{
				ServerAddr: flags.NetAddress{
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
				"-r=30",
			},
			wantFlags: wantFlags{
				ServerAddr:     standardFlags.ServerAddr,
				ReportInterval: 30,
				PollInterval:   standardFlags.PollInterval,
				RateLimit:      standardFlags.RateLimit,
				Key:            standardFlags.Key,
			},
		},
		{
			name: "checking non-standard polling interval flag",
			args: []string{
				"-p=2",
			},
			wantFlags: wantFlags{
				ServerAddr:     standardFlags.ServerAddr,
				ReportInterval: standardFlags.ReportInterval,
				PollInterval:   2,
				RateLimit:      standardFlags.RateLimit,
				Key:            standardFlags.Key,
			},
		},
		{
			name: "checking non-standard rate limit flag",
			args: []string{
				"-l=3",
			},
			wantFlags: wantFlags{
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
			wantFlags: wantFlags{
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
				"-r=5",
				"-p=1",
				"-k=non-standard-key",
			},
			wantFlags: wantFlags{
				ServerAddr: flags.NetAddress{
					Host: "example.com",
					Port: 80,
				},
				ReportInterval: 5,
				PollInterval:   1,
				RateLimit:      10,
				Key:            "non-standard-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseFlags("program", tt.args)

			if !assert.ObjectsAreEqual(tt.wantFlags.ServerAddr, flagServerAddr) {
				t.Errorf("expected flagServerAddr: %v, got: %v", tt.wantFlags.ServerAddr, flagServerAddr)
			}
			assert.Equalf(t, tt.wantFlags.ReportInterval, flagReportInterval, `expected flagReportInterval: %v, got: %v`, tt.wantFlags.ReportInterval, flagReportInterval)
			assert.Equalf(t, tt.wantFlags.PollInterval, flagPollInterval, `expected flagPollInterval: %v, got: %v`, tt.wantFlags.PollInterval, flagPollInterval)
			assert.Equalf(t, tt.wantFlags.RateLimit, flagRateLimit, `expected flagRateLimit: %v, got: %v`, tt.wantFlags.RateLimit, flagRateLimit)
			assert.Equalf(t, tt.wantFlags.Key, flagKey, `expected flagKey: "%v", got: "%v"`, tt.wantFlags.Key, flagKey)
		})
	}
}
