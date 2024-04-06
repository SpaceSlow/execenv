package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/SpaceSlow/execenv/cmd/storages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricRouter(t *testing.T) {
	
	type fields struct {
		method  string
		storage storages.MetricStorage
		path    string
	}
	type want struct {
		statusCode int
		metric     *metrics.Metric
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "incorrect metric type",
			fields: fields{
				method: http.MethodPost,
				path:   "/update/unknown_type/metric/42",
			},
			want: want{
				http.StatusBadRequest,
				nil,
			},
		},
		{
			name: "empty metric type",
			fields: fields{
				method: http.MethodPost,
				path:   "/update/",
			},
			want: want{
				http.StatusBadRequest,
				nil,
			},
		},
		{
			name: "root path",
			fields: fields{
				method: http.MethodPost,
				path:   "/",
			},
			want: want{
				http.StatusBadRequest,
				nil,
			},
		},
		{
			name: "incorrect method",
			fields: fields{
				method:  http.MethodGet,
				storage: storages.NewMemStorage(),
				path:    "/update/gauge/metric/1.21",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				metric:     nil,
			},
		},
		{
			name: "empty metric",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/gauge",
			},
			want: want{
				statusCode: http.StatusNotFound,
				metric:     nil,
			},
		},
		{
			name: "empty value",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/gauge/metricName/",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				metric:     nil,
			},
		},
		{
			name: "incorrect value (float for Counter metric type)",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/counter/metric/54.54",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				metric:     nil,
			},
		},
		{
			name: "incorrect value (string for Counter metric type)",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/counter/metric/incorrect_value",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				metric:     nil,
			},
		},
		{
			name: "incorrect value (string for Gauge metric type)",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/gauge/metric/incorrect_value",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				metric:     nil,
			},
		},
		{
			name: "correct value (Counter metric type)",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/counter/metric/101109",
			},
			want: want{
				statusCode: http.StatusOK,
				metric: &metrics.Metric{
					Type:  metrics.Counter,
					Name:  "metric",
					Value: int64(101109),
				},
			},
		},
		{
			name: "correct value (Gauge metric type)",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/gauge/metric/1.21",
			},
			want: want{
				statusCode: http.StatusOK,
				metric: &metrics.Metric{
					Type:  metrics.Gauge,
					Name:  "metric",
					Value: 1.21,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(MetricRouter(test.fields.storage))
			defer ts.Close()

			req, err := http.NewRequest(test.fields.method, ts.URL+test.fields.path, nil)
			require.NoError(t, err)

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, res.StatusCode)

			if res.StatusCode != http.StatusOK {
				return
			}

			storagedMetric, ok := test.fields.storage.Get(test.want.metric.Type, test.want.metric.Name)

			require.Equal(t, true, ok)

			assert.Equal(t, storagedMetric.Type, test.want.metric.Type)
			assert.Equal(t, storagedMetric.Name, test.want.metric.Name)
			assert.Equal(t, storagedMetric.Value, test.want.metric.Value)
		})
	}
}
