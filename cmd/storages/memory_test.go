package storages

import (
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
