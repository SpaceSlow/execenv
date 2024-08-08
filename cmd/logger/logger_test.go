package logger

import "testing"

func TestInitialize(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{
			name:    "logger with info level",
			level:   "info",
			wantErr: false,
		},
		{
			name:    "logger with debug level",
			level:   "debug",
			wantErr: false,
		},
		{
			name:    "logger with unknown level",
			level:   "unknown",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Initialize(tt.level); (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
