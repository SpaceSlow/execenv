package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
