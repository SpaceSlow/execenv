package metrics

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAddress[T int64 | float64](v T) *T {
	return &v
}

func TestSendMetrics(t *testing.T) {
	tests := []struct {
		name             string
		metrics          []Metric
		wantResponseBody []JSONMetric
	}{
		{
			name:             "empty metrics",
			metrics:          nil,
			wantResponseBody: []JSONMetric{},
		},
		{
			name: "one metric",
			metrics: []Metric{{
				Type:  Gauge,
				Name:  "RandomValue",
				Value: 1.21,
			}},
			wantResponseBody: []JSONMetric{
				{
					ID:    "RandomValue",
					MType: "gauge",
					Value: getAddress(1.21),
				},
			},
		},
		{
			name: "more metrics",
			metrics: []Metric{
				{
					Type:  Gauge,
					Name:  "RandomValue",
					Value: 1.21,
				},
				{
					Type:  Counter,
					Name:  "PollCount",
					Value: int64(42),
				},
				{
					Type:  Gauge,
					Name:  "SecondRandomValue",
					Value: 1011.09,
				},
			},
			wantResponseBody: []JSONMetric{
				{
					ID:    "RandomValue",
					MType: "gauge",
					Value: getAddress(1.21),
				},
				{
					ID:    "PollCount",
					MType: "counter",
					Delta: getAddress(int64(42)),
				},
				{
					ID:    "SecondRandomValue",
					MType: "gauge",
					Value: getAddress(1011.09),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			receivedMetrics := make([]JSONMetric, 0, len(test.metrics))
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, r.Method, http.MethodPost)
				require.Equal(t, r.Header.Get("Content-Type"), "application/json")
				require.Equal(t, r.Header.Get("Content-Encoding"), "gzip")

				dBody, err := gzip.NewReader(r.Body)
				require.NoError(t, err)

				var m []JSONMetric
				require.NoError(t, json.NewDecoder(dBody).Decode(&m))

				receivedMetrics = append(receivedMetrics, m...)
			}))
			defer testServer.Close()

			require.NoError(t, SendMetrics(testServer.URL+"/update/", test.metrics))

			assert.ElementsMatch(t, receivedMetrics, test.wantResponseBody)
		})
	}
}
