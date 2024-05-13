package storages

import (
	"sync"
	"testing"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_Add(t *testing.T) {
	type fields struct {
		counters map[string]int64
		gauges   map[string]float64
		metric   metrics.Metric
	}
	type want struct {
		err   bool
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "adding Counter metric in empty storage",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   make(map[string]float64),
				metric: metrics.Metric{
					Type:  metrics.Counter,
					Name:  "PollCount",
					Value: int64(5),
				},
			},
			want: want{
				err:   false,
				value: int64(5),
			},
		},
		{
			name: "adding Gauge metric in empty storage",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   make(map[string]float64),
				metric: metrics.Metric{
					Type:  metrics.Gauge,
					Name:  "RandomValue",
					Value: 1.21,
				},
			},
			want: want{
				err:   false,
				value: 1.21,
			},
		},
		{
			name: "summation of Counter metrics value",
			fields: fields{
				counters: map[string]int64{"PollCount": 10},
				gauges:   make(map[string]float64),
				metric: metrics.Metric{
					Type:  metrics.Counter,
					Name:  "PollCount",
					Value: int64(5),
				},
			},
			want: want{
				err:   false,
				value: int64(15),
			},
		},
		{
			name: "substitution of Gauge metrics value",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   map[string]float64{"RandomValue": 1.2},
				metric: metrics.Metric{
					Type:  metrics.Gauge,
					Name:  "PollCount",
					Value: 1.21,
				},
			},
			want: want{
				err:   false,
				value: 1.21,
			},
		},
		{
			name: "adding incorrect type metric",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   make(map[string]float64),
				metric: metrics.Metric{
					Type:  metrics.MetricType(-1),
					Name:  "PollCount",
					Value: int64(5),
				},
			},
			want: want{
				err:   true,
				value: nil,
			},
		},
		{
			name: "adding Counter metric with incorrect value",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   make(map[string]float64),
				metric: metrics.Metric{
					Type:  metrics.Counter,
					Name:  "PollCount",
					Value: "incorrect value",
				},
			},
			want: want{
				err:   true,
				value: nil,
			},
		},
		{
			name: "adding Gauge metric with incorrect value",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   make(map[string]float64),
				metric: metrics.Metric{
					Type:  metrics.Gauge,
					Name:  "PollCount",
					Value: "incorrect value",
				},
			},
			want: want{
				err:   true,
				value: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewMemStorage()
			storage.counters = test.fields.counters
			storage.gauges = test.fields.gauges

			_, err := storage.Add(&test.fields.metric)
			assert.Equal(t, err != nil, test.want.err)

			if err != nil {
				return
			}

			switch test.fields.metric.Type {
			case metrics.Counter:
				value, ok := storage.counters[test.fields.metric.Name]
				require.True(t, ok)
				assert.Equal(t, test.want.value, value)
			case metrics.Gauge:
				value, ok := storage.gauges[test.fields.metric.Name]
				require.True(t, ok)
				assert.Equal(t, test.want.value, value)
			}
		})
	}
}

func TestMemStorage_Get(t *testing.T) {
	type fields struct {
		mu       sync.Mutex
		counters map[string]int64
		gauges   map[string]float64
	}
	type args struct {
		metricType metrics.MetricType
		name       string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedMetric *metrics.Metric
		expectedOk     bool
	}{
		{
			name: "get existing counter metric",
			fields: fields{
				counters: map[string]int64{"PollCount": 10},
				gauges:   make(map[string]float64),
			},
			args: args{
				metricType: metrics.Counter,
				name:       "PollCount",
			},
			expectedMetric: &metrics.Metric{
				Type:  metrics.Counter,
				Name:  "PollCount",
				Value: int64(10),
			},
			expectedOk: true,
		},
		{
			name: "get existing gauge metric",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   map[string]float64{"RandomValue": 11.11},
			},
			args: args{
				metricType: metrics.Gauge,
				name:       "RandomValue",
			},
			expectedMetric: &metrics.Metric{
				Type:  metrics.Gauge,
				Name:  "RandomValue",
				Value: float64(11.11),
			},
			expectedOk: true,
		},
		{
			name: "get not-existing counter metric",
			fields: fields{
				counters: map[string]int64{"PollCount": 10},
				gauges:   make(map[string]float64),
			},
			args: args{
				metricType: metrics.Counter,
				name:       "NotExistMetric",
			},
			expectedMetric: nil,
			expectedOk:     false,
		},
		{
			name: "get not-existing metric for gauge",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   map[string]float64{"RandomValue": 11.11},
			},
			args: args{
				metricType: metrics.Gauge,
				name:       "NotExistMetric",
			},
			expectedMetric: nil,
			expectedOk:     false,
		},
		{
			name: "get incorrect type metric for counter type",
			fields: fields{
				counters: map[string]int64{"PollCount": 10},
				gauges:   make(map[string]float64),
			},
			args: args{
				metricType: metrics.MetricType(-1),
				name:       "PollCount",
			},
			expectedMetric: nil,
			expectedOk:     false,
		},
		{
			name: "get incorrect type metric for gauge type",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   map[string]float64{"RandomValue": 11.11},
			},
			args: args{
				metricType: metrics.MetricType(-1),
				name:       "RandomValue",
			},
			expectedMetric: nil,
			expectedOk:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemStorage()
			storage.counters = tt.fields.counters
			storage.gauges = tt.fields.gauges

			metric, ok := storage.Get(tt.args.metricType, tt.args.name)
			assert.Equal(t, tt.expectedMetric, metric)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}

func TestMemStorage_List(t *testing.T) {
	type fields struct {
		counters map[string]int64
		gauges   map[string]float64
	}
	tests := []struct {
		name        string
		fields      fields
		wantMetrics []metrics.Metric
	}{
		{
			name: "list only gauge metrics",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   map[string]float64{"RandomValue": float64(11.11), "RandomValue2": float64(66.66)},
			},
			wantMetrics: []metrics.Metric{
				{Type: metrics.Gauge, Name: "RandomValue2", Value: float64(66.66)},
				{Type: metrics.Gauge, Name: "RandomValue", Value: float64(11.11)},
			},
		},
		{
			name: "list only counter metrics",
			fields: fields{
				counters: map[string]int64{"PollCount": int64(7), "ExistingCounter": int64(9)},
				gauges:   make(map[string]float64),
			},
			wantMetrics: []metrics.Metric{
				{Type: metrics.Counter, Name: "PollCount", Value: int64(7)},
				{Type: metrics.Counter, Name: "ExistingCounter", Value: int64(9)},
			},
		},
		{
			name: "list some metrics",
			fields: fields{
				counters: map[string]int64{"PollCount": int64(7), "ExistingCounter": int64(9)},
				gauges:   map[string]float64{"RandomValue": float64(11.11), "RandomValue2": float64(66.66)},
			},
			wantMetrics: []metrics.Metric{
				{Type: metrics.Counter, Name: "PollCount", Value: int64(7)},
				{Type: metrics.Gauge, Name: "RandomValue2", Value: float64(66.66)},
				{Type: metrics.Gauge, Name: "RandomValue", Value: float64(11.11)},
				{Type: metrics.Counter, Name: "ExistingCounter", Value: int64(9)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemStorage()
			storage.counters = tt.fields.counters
			storage.gauges = tt.fields.gauges

			assert.ElementsMatch(t, tt.wantMetrics, storage.List())
		})
	}
}

func TestMemStorage_Batch(t *testing.T) {
	type fields struct {
		counters map[string]int64
		gauges   map[string]float64
	}
	type args struct {
		metricSlice []metrics.Metric
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMetrics []metrics.Metric
		wantErr     error
	}{
		{
			name: "check update values after batch one gauge metric",
			fields: fields{
				counters: make(map[string]int64),
				gauges:   map[string]float64{"RandomValue": float64(11.11)},
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
				counters: map[string]int64{"PollCount": 5},
				gauges:   make(map[string]float64),
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
				counters: map[string]int64{"PollCount": int64(5), "CounterMetric": int64(9)},
				gauges:   map[string]float64{"RandomValue": float64(65), "RandomValue2": float64(10.001)},
			},
			args: args{
				metricSlice: []metrics.Metric{
					{Type: metrics.Counter, Name: "PollCount", Value: int64(50)},
					{Type: metrics.Gauge, Name: "RandomValue2", Value: float64(0.999)},
				},
			},
			wantMetrics: []metrics.Metric{
				{Type: metrics.Counter, Name: "PollCount", Value: int64(55)},
				{Type: metrics.Gauge, Name: "RandomValue2", Value: float64(0.999)},
			},
			wantErr: nil,
		},
		{
			name: "check batch with incorrect metrics",
			fields: fields{
				counters: map[string]int64{"PollCount": int64(5), "CounterMetric": int64(9)},
				gauges:   map[string]float64{"RandomValue": float64(65), "RandomValue2": float64(10.001)},
			},
			args: args{
				metricSlice: []metrics.Metric{
					{Type: metrics.MetricType(-1), Name: "PollCount", Value: int64(50)},
					{Type: metrics.Gauge, Name: "RandomValue2", Value: float64(0.999)},
				},
			},
			wantMetrics: nil,
			wantErr:     metrics.ErrIncorrectMetricTypeOrValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemStorage()
			storage.counters = tt.fields.counters
			storage.gauges = tt.fields.gauges

			metricSlice, err := storage.Batch(tt.args.metricSlice)
			require.Equal(t, tt.wantErr, err)
			assert.ElementsMatch(t, tt.wantMetrics, metricSlice)
		})
	}
}
