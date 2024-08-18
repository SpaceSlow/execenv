package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDuration_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name             string
		expectedDuration *Duration
		durationJson     []byte
		wantErr          bool
	}{
		{
			name:             "1 second",
			durationJson:     []byte(`"1s"`),
			expectedDuration: &Duration{time.Second},
			wantErr:          false,
		},
		{
			name:             "14 seconds 55 milliseconds",
			durationJson:     []byte(`"14s55ms"`),
			expectedDuration: &Duration{14*time.Second + 55*time.Millisecond},
			wantErr:          false,
		},
		{
			name:         "incorrect duration",
			durationJson: []byte(`"incorrect"`),
			wantErr:      true,
		},
		{
			name:         "incorrect type duration",
			durationJson: []byte(`{}`),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := new(Duration)

			err := d.UnmarshalJSON(tt.durationJson)
			require.Equal(t, tt.wantErr, err != nil)

			if err == nil {
				assert.Equal(t, d, tt.expectedDuration)
			}
		})
	}
}
