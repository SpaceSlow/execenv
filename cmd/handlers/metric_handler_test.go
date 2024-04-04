package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/SpaceSlow/execenv/cmd/storages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		method     string
		metricType metrics.MetricType
		storage    storages.Storage
		url        string
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
			name: "incorrect method",
			fields: fields{
				method:     http.MethodGet,
				metricType: metrics.Gauge,
				storage:    storages.NewMemStorage(),
				url:        "/update/gauge/metric/1.21",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				metric:     nil,
			},
		},
		{
			name: "empty metric",
			fields: fields{
				method:     http.MethodPost,
				metricType: metrics.Gauge,
				storage:    storages.NewMemStorage(),
				url:        "/update/gauge",
			},
			want: want{
				statusCode: http.StatusNotFound,
				metric:     nil,
			},
		},
		{
			name: "empty value",
			fields: fields{
				method:     http.MethodPost,
				metricType: metrics.Gauge,
				storage:    storages.NewMemStorage(),
				url:        "/update/gauge/metricName/",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				metric:     nil,
			},
		},
		{
			name: "incorrect value (float for Counter metric type)",
			fields: fields{
				method:     http.MethodPost,
				metricType: metrics.Counter,
				storage:    storages.NewMemStorage(),
				url:        "/update/counter/metric/54.54",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				metric:     nil,
			},
		},
		{
			name: "incorrect value (string for Counter metric type)",
			fields: fields{
				method:     http.MethodPost,
				metricType: metrics.Counter,
				storage:    storages.NewMemStorage(),
				url:        "/update/counter/metric/incorrect_value",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				metric:     nil,
			},
		},
		{
			name: "incorrect value (string for Gauge metric type)",
			fields: fields{
				method:     http.MethodPost,
				metricType: metrics.Gauge,
				storage:    storages.NewMemStorage(),
				url:        "/update/gauge/metric/incorrect_value",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				metric:     nil,
			},
		},
		{
			name: "correct value (Counter metric type)",
			fields: fields{
				method:     http.MethodPost,
				metricType: metrics.Counter,
				storage:    storages.NewMemStorage(),
				url:        "/update/counter/metric/101109",
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
				method:     http.MethodPost,
				metricType: metrics.Gauge,
				storage:    storages.NewMemStorage(),
				url:        "/update/gauge/metric/1.21",
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
			request := httptest.NewRequest(test.fields.method, test.fields.url, nil)

			w := httptest.NewRecorder()
			http.StripPrefix(fmt.Sprintf("/update/%s/", test.fields.metricType.String()), MetricHandler{
				MetricType: test.fields.metricType,
				Storage:    test.fields.storage,
			}).ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)

			if res.StatusCode != http.StatusOK {
				return
			}

			storagedMetric, ok := test.fields.storage.Get(test.want.metric.Name)

			require.Equal(t, true, ok)

			assert.Equal(t, storagedMetric.Type, test.want.metric.Type)
			assert.Equal(t, storagedMetric.Name, test.want.metric.Name)
			assert.Equal(t, storagedMetric.Value, test.want.metric.Value)
		})
	}
}
