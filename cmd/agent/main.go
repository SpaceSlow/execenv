package main

import (
	"math/rand"
	"time"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

func main() {
	cfg, err := GetConfigWithFlags()

	if err != nil {
		panic(err)
	}

	var metricSlice []metrics.Metric

	url := "http://" + cfg.ServerAddr.String() + "/update/"
	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

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
