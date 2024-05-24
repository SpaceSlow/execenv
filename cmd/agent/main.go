package main

import (
	"log"
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
	metricSender := metrics.NewMetricSender(url, cfg.Key)
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
			metricSender.Push(metricSlice)
		case <-reportTick:
			select {
			case err := <-metrics.RetryFunc(metricSender.Send, cfg.Delays):
				if err != nil {
					log.Printf("sending error: %s", err)
					continue
				}
				log.Printf("sended metrics")
			}
			pollCount = 0
		}
	}
}
