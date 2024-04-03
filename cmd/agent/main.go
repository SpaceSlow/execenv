package main

import (
	"math/rand"
	"time"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	var metricSlice []metrics.Metric
	url := "http://localhost:8080/update"

	pollCount := 0

	for controlInterval := reportInterval; ; controlInterval -= pollInterval {
		if controlInterval <= time.Duration(0) {
			metrics.SendMetrics(url, metricSlice)
			pollCount = 0
			controlInterval = reportInterval
		}

		pollCount++
		metricSlice = metrics.GetRuntimeMetrics()
		metricSlice = append(
			metricSlice,
			metrics.Metric{
				Type:  metrics.Gauge,
				Name:  "RandomValue",
				Value: rand.Float64(),
			},
			metrics.Metric{
				Type:  metrics.Counter,
				Name:  "PollCount",
				Value: pollCount,
			},
		)

		time.Sleep(pollInterval)
	}
}
