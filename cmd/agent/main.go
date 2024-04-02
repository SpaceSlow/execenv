package main

import (
	"sync"
	"time"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

const (
	pollInterval   = time.Duration(2)
	reportInterval = time.Duration(10)
)

func main() {
	var metricSlice []metrics.Metric
	var mu sync.Mutex

	pollTicker := time.NewTicker(time.Second * pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Second * reportInterval)
	defer reportTicker.Stop()
	for {
		select {
		case <-pollTicker.C:
			mu.Lock()
			metricSlice = metrics.GetMetrics()
			mu.Unlock()
		case <-reportTicker.C:
			mu.Lock()
			metrics.SendMetrics(metricSlice)
			mu.Unlock()
		}
	}
}
