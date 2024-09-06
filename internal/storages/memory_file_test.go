package storages

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SpaceSlow/execenv/internal/metrics"
)

func TestMemFileStorage_Add(t *testing.T) {
	tests := []struct {
		wantErr error
		metric  *metrics.Metric
		name    string
	}{
		{
			name: "check adding correct metric",
			metric: &metrics.Metric{
				Type:  metrics.Gauge,
				Name:  "RandomValue",
				Value: 7.07,
			},
			wantErr: nil,
		},
		{
			name: "check adding metric with incorrect type",
			metric: &metrics.Metric{
				Type:  -1,
				Name:  "RandomValue2",
				Value: 7.07,
			},
			wantErr: metrics.ErrIncorrectMetricTypeOrValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewMemFileStorage(context.Background(), path.Join(os.TempDir(), randStringBytes(10)), 0, false)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, s.Close())
				require.NoError(t, os.Remove(s.f.Name()))
			}()

			data := make([]byte, 100)
			size, err := s.f.Read(data)
			require.ErrorIs(t, err, io.EOF)
			assert.Zero(t, size)

			_, err = s.Add(tt.metric)
			require.ErrorIs(t, err, tt.wantErr)
			if err != nil {
				return
			}

			size, err = s.f.Read(data)
			require.NoError(t, err)
			assert.Greater(t, size, 0)
		})
	}
}

func TestMemFileStorage_Batch(t *testing.T) {
	type fields struct {
		counters counters
		gauges   gauges
	}
	type args struct {
		metricSlice []metrics.Metric
	}
	tests := []struct {
		wantErr     error
		name        string
		fields      fields
		args        args
		wantMetrics []metrics.Metric
	}{
		{
			name: "check update values after batch one gauge metric",
			fields: fields{
				counters: make(counters),
				gauges:   gauges{"RandomValue": 11.11},
			},
			args: args{
				metricSlice: []metrics.Metric{{Type: metrics.Gauge, Name: "RandomValue", Value: float64(55.55)}},
			},
			wantMetrics: []metrics.Metric{{Type: metrics.Gauge, Name: "RandomValue", Value: float64(55.55)}},
			wantErr:     nil,
		},
		{
			name: "check summing deltas after batch one counter metric",
			fields: fields{
				counters: counters{"PollCount": 5},
				gauges:   make(gauges),
			},
			args: args{
				metricSlice: []metrics.Metric{{Type: metrics.Counter, Name: "PollCount", Value: int64(50)}},
			},
			wantMetrics: []metrics.Metric{{Type: metrics.Counter, Name: "PollCount", Value: int64(55)}},
			wantErr:     nil,
		},
		{
			name: "check affect on existing metrics after batch more metrics",
			fields: fields{
				counters: counters{"PollCount": 5, "CounterMetric": 9},
				gauges:   gauges{"RandomValue": 65, "RandomValue2": 10.001},
			},
			args: args{
				metricSlice: []metrics.Metric{
					{Type: metrics.Counter, Name: "PollCount", Value: int64(50)},
					{Type: metrics.Gauge, Name: "RandomValue2", Value: float64(0.999)},
				},
			},
			wantMetrics: []metrics.Metric{
				{Type: metrics.Counter, Name: "PollCount", Value: int64(55)},
				{Type: metrics.Counter, Name: "CounterMetric", Value: int64(9)},
				{Type: metrics.Gauge, Name: "RandomValue", Value: float64(65)},
				{Type: metrics.Gauge, Name: "RandomValue2", Value: float64(0.999)},
			},
			wantErr: nil,
		},
		{
			name: "check batch with incorrect metrics",
			fields: fields{
				counters: counters{"PollCount": 5, "CounterMetric": 9},
				gauges:   gauges{"RandomValue": 65, "RandomValue2": 10.001},
			},
			args: args{
				metricSlice: []metrics.Metric{
					{Type: metrics.MetricType(-1), Name: "PollCount", Value: int64(50)},
					{Type: metrics.Gauge, Name: "RandomValue2", Value: float64(0.999)},
				},
			},
			wantMetrics: []metrics.Metric{
				{Type: metrics.Counter, Name: "PollCount", Value: int64(5)},
				{Type: metrics.Counter, Name: "CounterMetric", Value: int64(9)},
				{Type: metrics.Gauge, Name: "RandomValue", Value: float64(65)},
				{Type: metrics.Gauge, Name: "RandomValue2", Value: float64(10.001)},
			},
			wantErr: metrics.ErrIncorrectMetricTypeOrValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewMemFileStorage(context.Background(), path.Join(os.TempDir(), randStringBytes(10)), 0, false)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, s.Close())
				require.NoError(t, os.Remove(s.f.Name()))
			}()

			s.counters = tt.fields.counters
			s.gauges = tt.fields.gauges
			require.NoError(t, s.SaveMetricsToFile())

			err = s.Batch(tt.args.metricSlice)
			require.Equal(t, tt.wantErr, err)

			data, err := os.ReadFile(s.f.Name())
			require.NoError(t, err)

			gotMetrics := make([]metrics.Metric, 0)
			require.NoError(t, json.Unmarshal(data, &gotMetrics))
			assert.ElementsMatch(t, tt.wantMetrics, gotMetrics)
		})
	}
}

func TestMemFileStorage_Close(t *testing.T) {
	s, err := NewMemFileStorage(context.Background(), path.Join(os.TempDir(), randStringBytes(10)), 0, false)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(s.f.Name()))
	}()

	assert.NoError(t, s.Close())
}

func TestMemFileStorage_LoadMetricsFromFile(t *testing.T) {
	type wantMetrics struct {
		counters counters
		gauges   gauges
	}
	tests := []struct {
		name        string
		wantMetrics wantMetrics
		wantErr     error
		data        []byte
	}{
		{
			name:        "loading one metric from file",
			data:        []byte(`[{"id": "RandomValue", "type": "gauge", "value": 7.07}]`),
			wantMetrics: wantMetrics{counters: make(counters), gauges: gauges{"RandomValue": 7.07}},
			wantErr:     nil,
		},
		{
			name: "loading several metrics from file",
			data: []byte(`[
			   {"id": "RandomValue", "type": "gauge", "value": 7.07},
			   {"id": "PollCount", "type": "counter", "delta": 10},
			   {"id": "RandomValue-2", "type": "gauge", "value": 9.37}
			]`),
			wantMetrics: wantMetrics{counters: counters{"PollCount": 10}, gauges: gauges{"RandomValue": 7.07, "RandomValue-2": 9.37}},
			wantErr:     nil,
		},
		{
			name:        "error occurred when loading incorrect metric from file",
			data:        []byte(`[{"id": "RandomValue", "type": "incorrect", "value": 7.07}]`),
			wantMetrics: wantMetrics{counters: make(counters), gauges: make(gauges)},
			wantErr:     metrics.ErrIncorrectMetricTypeOrValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewMemFileStorage(context.Background(), path.Join(os.TempDir(), randStringBytes(10)), 0, false)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, s.Close())
				require.NoError(t, os.Remove(s.f.Name()))
			}()

			_, err = s.f.Write(tt.data)
			require.NoError(t, err)

			err = s.LoadMetricsFromFile()
			assert.ErrorIs(t, tt.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tt.wantMetrics.counters, s.counters)
			assert.Equal(t, tt.wantMetrics.gauges, s.gauges)
		})
	}
}

func TestMemFileStorage_SaveMetricsToFile(t *testing.T) {
	type fields struct {
		counters counters
		gauges   gauges
	}
	tests := []struct {
		fields      fields
		wantErr     error
		name        string
		wantMetrics string
	}{
		{
			name: "check saving metrics to file",
			fields: fields{
				counters: counters{"PollCount-42": 42},
				gauges:   gauges{"RandomValue-101": 1.01, "FictiveGauge-42": 4.2},
			},
			wantMetrics: `[
			   {"id": "PollCount-42", "type": "counter", "delta": 42},
			   {"id": "RandomValue-101", "type": "gauge", "value": 1.01},
			   {"id": "FictiveGauge-42", "type": "gauge", "value": 4.2}
			]`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewMemFileStorage(context.Background(), path.Join(os.TempDir(), randStringBytes(10)), 1000*time.Second, false)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, s.Close())
				require.NoError(t, os.Remove(s.f.Name()))
			}()

			s.counters = tt.fields.counters
			s.gauges = tt.fields.gauges

			err = s.SaveMetricsToFile()
			require.NoError(t, err)

			data, err := os.ReadFile(s.f.Name())
			require.NoError(t, err)

			var expectedMetrics []metrics.Metric
			err = json.Unmarshal([]byte(tt.wantMetrics), &expectedMetrics)
			require.NoError(t, err)

			var actualMetrics []metrics.Metric
			err = json.Unmarshal(data, &actualMetrics)
			require.NoError(t, err)

			assert.ElementsMatch(t, expectedMetrics, actualMetrics)
		})
	}
}
