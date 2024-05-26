package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
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
	metricsCh := make(chan []metrics.Metric) // TODO: goroutines for get metrics and send metrics
	closed := make(chan os.Signal, 1)
	defer close(closed)
	signal.Notify(closed, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	for {
		select {
		case <-pollTick:
			go metrics.GetMetrics(metricsCh)
			metricSlice = <-metricsCh
			pollCount++
			metricSlice = append(metricSlice, metrics.Metric{
				Type:  metrics.Counter,
				Name:  "PollCount",
				Value: pollCount,
			})
			metricSender.Push(metricSlice)
		case <-reportTick:
			err := <-metrics.RetryFunc(metricSender.Send, cfg.Delays)
			if err != nil {
				log.Printf("sending error: %s", err)
				continue
			}
			log.Printf("sended metrics")
			pollCount = 0
		case <-closed:
			close(metricsCh)
			log.Printf("stopped agent")
		}
	}
}
