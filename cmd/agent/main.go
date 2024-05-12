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

	url := "http://" + cfg.ServerAddr.String() + "/updates/"
	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	pollTick := time.Tick(pollInterval)
	reportTick := time.Tick(reportInterval)
	var pollCount int64

	for {
		select {
		case <-pollTick:
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
		case <-reportTick:
			err := metrics.SendMetrics(url, metricSlice)
			if err != nil {
				go retrySendMetrics(url, metricSlice)
			}
			pollCount = 0
		}
	}
}

func retrySendMetrics(url string, metricSlice []metrics.Metric) error {
	var err error
	delays := []time.Duration{
		time.Second,
		3 * time.Second,
		5 * time.Second,
	}
	for attempt := 0; attempt < len(delays); attempt++ {
		time.Sleep(delays[attempt])
		err = metrics.SendMetrics(url, metricSlice)
		if err != nil {
			attempt++
		} else {
			return err
		}
	}
	return err
}
