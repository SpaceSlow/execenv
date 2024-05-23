package main

import (
	"math/rand"
	"time"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

var pollCount int64

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
	//metricsCh := make(chan []metrics.Metric) TODO: goroutines for get metrics and send metrics

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
			err := metrics.SendMetrics(url, cfg.Key, metricSlice)
			if err != nil {
				go retrySendMetrics(cfg, url, metricSlice)
			}
			pollCount = 0
		}
	}
}

func retrySendMetrics(cfg *Config, url string, metricSlice []metrics.Metric) error {
	var err error
	for attempt := 0; attempt < len(cfg.Delays); attempt++ {
		time.Sleep(cfg.Delays[attempt])
		if err := metrics.SendMetrics(url, cfg.Key, metricSlice); err == nil {
			return nil
		}
		attempt++
	}
	return err
}
