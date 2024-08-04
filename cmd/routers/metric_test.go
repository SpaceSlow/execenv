package routers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

type fields struct {
	method  string
	storage storages.MetricStorage
	path    string
	body    string
}

type want struct {
	body       string
	bodyJSON   string
	statusCode int
}

type testCase struct {
	name   string
	fields fields
	want   want
}

func newMemStorageWithMetrics(metrics []metrics.Metric) *storages.MemStorage {
	storage := storages.NewMemStorage()

	for _, metric := range metrics {
		storage.Add(&metric)
	}

	return storage
}

func TestMetricRouter(t *testing.T) {
	tests := []testCase{
		{
			name: "incorrect metric type",
			fields: fields{
				method: http.MethodPost,
				path:   "/update/unknown_type/metric/42",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "empty metric type",
			fields: fields{
				method: http.MethodPost,
				path:   "/update/",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "root path",
			fields: fields{
				method: http.MethodPost,
				path:   "/",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
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
			},
		},
		{
			name: "get existing metric (Gauge metric type)",
			fields: fields{
				method: http.MethodGet,
				storage: newMemStorageWithMetrics([]metrics.Metric{
					{
						Type:  metrics.Gauge,
						Name:  "RandomValue",
						Value: 3.01,
					},
				}),
				path: "/value/gauge/RandomValue",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "get existing metric (Counter metric type)",
			fields: fields{
				method: http.MethodGet,
				storage: newMemStorageWithMetrics([]metrics.Metric{
					{
						Type:  metrics.Counter,
						Name:  "PollCount",
						Value: int64(15),
					},
				}),
				path: "/value/counter/PollCount",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "get non-existing metric (Counter metric type)",
			fields: fields{
				method: http.MethodGet,
				storage: newMemStorageWithMetrics([]metrics.Metric{
					{
						Type:  metrics.Gauge,
						Name:  "RandomValue",
						Value: 3.01,
					},
				}),
				path: "/value/counter/RandomValue",
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "get non-existing metric (Gauge metric type)",
			fields: fields{
				method: http.MethodGet,
				storage: newMemStorageWithMetrics([]metrics.Metric{
					{
						Type:  metrics.Gauge,
						Name:  "RandomValue",
						Value: 2.97,
					},
				}),
				path: "/value/counter/ValueRandom",
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "getting all existing metrics",
			fields: fields{
				method: http.MethodGet,
				storage: newMemStorageWithMetrics([]metrics.Metric{
					{
						Type:  metrics.Gauge,
						Name:  "RandomValue",
						Value: 2.97,
					},
					{
						Type:  metrics.Counter,
						Name:  "PollCount",
						Value: int64(10),
					},
				}),
				path: "/",
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "PollCount = 10 (counter)\nRandomValue = 2.97 (gauge)\n",
			},
		},
		{
			name: "update counter metric via json request",
			fields: fields{
				method: http.MethodPost,
				storage: newMemStorageWithMetrics([]metrics.Metric{
					{
						Type:  metrics.Counter,
						Name:  "PollCount",
						Value: int64(10),
					},
				}),
				path: "/update/",
				body: `{"id": "PollCount", "type": "counter", "delta": 5}`,
			},
			want: want{
				statusCode: http.StatusOK,
				bodyJSON:   `{"id":"PollCount","type":"counter","delta":15}`,
			},
		},
		{
			name: "update incorrect metric type via json request",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/",
				body:    `{"id": "PollCount", "type": "incorrect_type", "delta": 5}`,
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "incorrect filled fields in json request (Delta in Gauge)",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/",
				body:    `{"id": "RandomValue", "type": "gauge", "delta": 5}`,
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "incorrect filled fields in json request (Value in Counter)",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/",
				body:    `{"id": "PollCount", "type": "counter", "value": 117.9}`,
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "get not-existing metric via json request",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/value/",
				body:    `{"id": "PollCount", "type": "counter"}`,
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "get existing metric via json request",
			fields: fields{
				method: http.MethodPost,
				storage: newMemStorageWithMetrics([]metrics.Metric{
					{
						Type:  metrics.Counter,
						Name:  "PollCount",
						Value: int64(10),
					},
				}),
				path: "/value/",
				body: `{"id": "PollCount", "type": "counter"}`,
			},
			want: want{
				statusCode: http.StatusOK,
				bodyJSON:   `{"id":"PollCount","type":"counter","delta":10}`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(MetricRouter(test.fields.storage))
			defer ts.Close()

			req, err := http.NewRequest(test.fields.method, ts.URL+test.fields.path, nil)
			require.NoError(t, err)

			if len(test.fields.body) > 0 {
				req.Header.Set("Content-Type", "application/json")
				req.Body = io.NopCloser(strings.NewReader(test.fields.body))
			}

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, res.StatusCode)

			if test.want.body != "" {
				body, err := io.ReadAll(res.Body)
				require.NoError(t, res.Body.Close())
				require.NoError(t, err)
				if test.want.bodyJSON != "" {
					assert.JSONEq(t, test.want.bodyJSON, string(body))
				} else {
					assert.Equal(t, test.want.body, string(body))
				}
			}
		})
	}
}
