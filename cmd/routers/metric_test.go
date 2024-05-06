package routers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/SpaceSlow/execenv/cmd/storages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fields struct {
	method     string
	storage    storages.MetricStorage
	path       string
	jsonMetric *metrics.JSONMetric
}

type want struct {
	statusCode int
	metric     *metrics.Metric
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

func requireEqualExistingMetricsWithResponse(t *testing.T, test testCase, res *http.Response) {
	body, err := io.ReadAll(res.Body)
	require.NoError(t, res.Body.Close())
	require.NoError(t, err)
	lines := strings.FieldsFunc(string(body), func(r rune) bool {
		return r == '\n'
	})

	metricStrings := make([]string, 0, len(lines))
	for _, metric := range test.fields.storage.List() {
		metricStrings = append(metricStrings, metric.String())
	}

	require.ElementsMatch(t, metricStrings, lines)
}

func getAddress[T int64 | float64](v T) *T {
	return &v
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
				metric: &metrics.Metric{
					Type:  metrics.Gauge,
					Name:  "RandomValue",
					Value: 3.01,
				},
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
				metric: &metrics.Metric{
					Type:  metrics.Counter,
					Name:  "PollCount",
					Value: int64(15),
				},
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
				jsonMetric: &metrics.JSONMetric{
					ID:    "PollCount",
					MType: "counter",
					Delta: getAddress(int64(5)),
				},
			},
			want: want{
				statusCode: http.StatusOK,
				metric: &metrics.Metric{
					Type:  metrics.Counter,
					Name:  "PollCount",
					Value: int64(15),
				},
			},
		},
		{
			name: "update incorrect metric type via json request",
			fields: fields{
				method:  http.MethodPost,
				storage: storages.NewMemStorage(),
				path:    "/update/",
				jsonMetric: &metrics.JSONMetric{
					ID:    "PollCount",
					MType: "incorrect_type",
					Delta: getAddress(int64(5)),
				},
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
				jsonMetric: &metrics.JSONMetric{
					ID:    "RandomValue",
					MType: "gauge",
					Delta: getAddress(int64(5)),
				},
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
				jsonMetric: &metrics.JSONMetric{
					ID:    "PollCount",
					MType: "counter",
					Value: getAddress(117.9),
				},
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
				jsonMetric: &metrics.JSONMetric{
					ID:    "PollCount",
					MType: "counter",
				},
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
				jsonMetric: &metrics.JSONMetric{
					ID:    "PollCount",
					MType: "counter",
				},
			},
			want: want{
				statusCode: http.StatusOK,
				metric: &metrics.Metric{
					Type:  metrics.Counter,
					Name:  "PollCount",
					Value: int64(10),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(MetricRouter(test.fields.storage))
			defer ts.Close()

			var body io.Reader

			isJSONRequest := (test.fields.path == "/update/" || test.fields.path == "/value/") && test.fields.jsonMetric != nil
			if isJSONRequest {
				m, err := json.Marshal(test.fields.jsonMetric)
				require.NoError(t, err)
				body = bytes.NewReader(m)
			}

			req, err := http.NewRequest(test.fields.method, ts.URL+test.fields.path, body)
			require.NoError(t, err)

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, res.StatusCode)

			if res.StatusCode != http.StatusOK {
				return
			}

			if test.fields.path == "/" {
				requireEqualExistingMetricsWithResponse(t, test, res)
				return
			}

			if isJSONRequest {
				var data []byte
				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				var m metrics.Metric
				require.NoError(t, m.UnmarshalJSON(data))

				assert.Equal(t, m.Type, test.want.metric.Type)
				assert.Equal(t, m.Name, test.want.metric.Name)
				assert.Equal(t, m.Value, test.want.metric.Value)
			}

			storagedMetric, ok := test.fields.storage.Get(test.want.metric.Type, test.want.metric.Name)

			require.Equal(t, true, ok)

			assert.Equal(t, storagedMetric.Type, test.want.metric.Type)
			assert.Equal(t, storagedMetric.Name, test.want.metric.Name)
			assert.Equal(t, storagedMetric.Value, test.want.metric.Value)
		})
	}
}
