package metrics

import "testing"

func TestParseMetricType(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    MetricType
		wantErr bool
	}{
		{name: "counter metric type", arg: "counter", want: Counter, wantErr: false},
		{name: "gauge metric type", arg: "gauge", want: Gauge, wantErr: false},
		{name: "incorrect metric type", arg: "incorrect", want: MetricType(-1), wantErr: true},
		{name: "must lowercase for metric type", arg: "CouNter", want: MetricType(-1), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMetricType(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMetricType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseMetricType() got = %v, want %v", got, tt.want)
			}
		})
	}
}
