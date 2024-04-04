package storages

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

func TestMemStorage_Add(t *testing.T) {
	type fields struct {
		metrics map[string]interface{}
		metric  metrics.Metric
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
				metrics: make(map[string]interface{}),
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
				metrics: make(map[string]interface{}),
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
				metrics: map[string]interface{}{
					"PollCount": int64(10),
				},
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
				metrics: map[string]interface{}{
					"RandomValue": 1.2,
				},
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
				metrics: make(map[string]interface{}),
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
				metrics: make(map[string]interface{}),
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
				metrics: make(map[string]interface{}),
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
			storage.metrics = test.fields.metrics

			err := storage.Add(&test.fields.metric)
			assert.Equal(t, err != nil, test.want.err)

			if err != nil {
				return
			}

			value, ok := storage.metrics[test.fields.metric.Name]
			
			require.True(t, ok)
			assert.Equal(t, test.want.value, value)
		})
	}
}
