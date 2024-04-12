package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMetrics(t *testing.T) {
	tests := []struct {
		name     string
		metrics  []Metric
		wantUrls []string
	}{
		{
			name:     "empty metrics",
			metrics:  nil,
			wantUrls: []string{},
		},
		{
			name: "one metric",
			metrics: []Metric{{
				Type:  Gauge,
				Name:  "RandomValue",
				Value: 1.21,
			}},
			wantUrls: []string{
				"/update/gauge/RandomValue/1.21",
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
					Value: 42,
				},
				{
					Type:  Gauge,
					Name:  "SecondRandomValue",
					Value: 1011.09,
				},
			},
			wantUrls: []string{
				"/update/gauge/RandomValue/1.21",
				"/update/counter/PollCount/42",
				"/update/gauge/SecondRandomValue/1011.09",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urls := make([]string, 0, len(test.metrics))
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				urls = append(urls, r.RequestURI)
			}))
			defer testServer.Close()

			SendMetrics(testServer.URL+"/update", test.metrics)

			assert.ElementsMatch(t, urls, test.wantUrls)
		})
	}
}
