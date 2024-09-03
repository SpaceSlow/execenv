package worker

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"math/rand"
	"runtime"
	"sync/atomic"

	"github.com/SpaceSlow/execenv/internal/client"
	"github.com/SpaceSlow/execenv/internal/config"
	"github.com/SpaceSlow/execenv/internal/metrics"
)

// MetricWorkers служит для аккумуляции и отправки метрик на сервер, с заданным ключом.
type MetricWorkers struct {
	errorsCh chan error

	client    *client.Client
	pollCount atomic.Int64
}

func NewMetricWorkers() (*MetricWorkers, error) {
	cfg, err := config.GetAgentConfig()
	if err != nil {
		return nil, err
	}
	client, err := client.NewClient()
	if err != nil {
		return nil, err
	}
	return &MetricWorkers{
		client:   client,
		errorsCh: make(chan error, cfg.RateLimit),
	}, nil
}

func (mw *MetricWorkers) Send(metrics []metrics.Metric) {
	pollCount := mw.pollCount.Load()

	err := mw.client.Send(metrics)
	if err != nil {
		mw.errorsCh <- err
		return
	}
	mw.pollCount.Add(-pollCount)
	mw.errorsCh <- nil
}

func (mw *MetricWorkers) Poll(pollCh chan []metrics.Metric) {
	metricSlice := make([]metrics.Metric, 0)
	runtimeMetricsCh := mw.getRuntimeMetrics()
	gopsutilMetricsCh := mw.getGopsutilMetrics()
	for m := range metrics.FanIn(runtimeMetricsCh, gopsutilMetricsCh) {
		metricSlice = append(metricSlice, m...)
	}
	metricSlice = append(
		metricSlice,
		metrics.Metric{
			Type:  metrics.Gauge,
			Name:  "RandomValue",
			Value: rand.Float64(),
		},
	)
	metricSlice = append(metricSlice, metrics.Metric{
		Type:  metrics.Counter,
		Name:  "PollCount",
		Value: mw.pollCount.Add(1),
	})

	if len(pollCh) > 0 {
		<-pollCh
	}
	pollCh <- metricSlice
}

func (mw *MetricWorkers) Close() {
	close(mw.errorsCh)
}

func (mw *MetricWorkers) Err() chan error {
	return mw.errorsCh
}

func (mw *MetricWorkers) getGopsutilMetrics() chan []metrics.Metric {
	metricsCh := make(chan []metrics.Metric)

	go func() {
		defer close(metricsCh)
		v, _ := mem.VirtualMemory()
		cpu, _ := cpu.Percent(0, false)

		metrics := []metrics.Metric{
			{
				Type:  metrics.Gauge,
				Name:  "TotalMemory",
				Value: float64(v.Total),
			},
			{
				Type:  metrics.Gauge,
				Name:  "FreeMemory",
				Value: float64(v.Free),
			},
			{
				Type:  metrics.Gauge,
				Name:  "CPUtilization1",
				Value: cpu[0],
			},
		}
		metricsCh <- metrics
	}()

	return metricsCh
}

func (mw *MetricWorkers) getRuntimeMetrics() chan []metrics.Metric {
	metricsCh := make(chan []metrics.Metric)

	go func() {
		defer close(metricsCh)
		var rtm runtime.MemStats
		runtime.ReadMemStats(&rtm)

		metrics := []metrics.Metric{
			{
				Type:  metrics.Gauge,
				Name:  "Alloc",
				Value: float64(rtm.Alloc),
			},
			{
				Type:  metrics.Gauge,
				Name:  "BuckHashSys",
				Value: float64(rtm.BuckHashSys),
			},
			{
				Type:  metrics.Gauge,
				Name:  "Frees",
				Value: float64(rtm.Frees),
			},
			{
				Type:  metrics.Gauge,
				Name:  "GCCPUFraction",
				Value: float64(rtm.GCCPUFraction),
			},
			{
				Type:  metrics.Gauge,
				Name:  "GCSys",
				Value: float64(rtm.GCSys),
			},
			{
				Type:  metrics.Gauge,
				Name:  "HeapAlloc",
				Value: float64(rtm.HeapAlloc),
			},
			{
				Type:  metrics.Gauge,
				Name:  "HeapIdle",
				Value: float64(rtm.HeapIdle),
			},
			{
				Type:  metrics.Gauge,
				Name:  "HeapInuse",
				Value: float64(rtm.HeapInuse),
			},
			{
				Type:  metrics.Gauge,
				Name:  "HeapObjects",
				Value: float64(rtm.HeapObjects),
			},
			{
				Type:  metrics.Gauge,
				Name:  "HeapReleased",
				Value: float64(rtm.HeapReleased),
			},
			{
				Type:  metrics.Gauge,
				Name:  "HeapSys",
				Value: float64(rtm.HeapSys),
			},
			{
				Type:  metrics.Gauge,
				Name:  "LastGC",
				Value: float64(rtm.LastGC),
			},
			{
				Type:  metrics.Gauge,
				Name:  "Lookups",
				Value: float64(rtm.Lookups),
			},
			{
				Type:  metrics.Gauge,
				Name:  "MCacheInuse",
				Value: float64(rtm.MCacheInuse),
			},
			{
				Type:  metrics.Gauge,
				Name:  "MCacheSys",
				Value: float64(rtm.MCacheSys),
			},
			{
				Type:  metrics.Gauge,
				Name:  "MSpanInuse",
				Value: float64(rtm.MSpanInuse),
			},
			{
				Type:  metrics.Gauge,
				Name:  "MSpanSys",
				Value: float64(rtm.MSpanSys),
			},
			{
				Type:  metrics.Gauge,
				Name:  "Mallocs",
				Value: float64(rtm.Mallocs),
			},
			{
				Type:  metrics.Gauge,
				Name:  "NextGC",
				Value: float64(rtm.NextGC),
			},
			{
				Type:  metrics.Gauge,
				Name:  "NumForcedGC",
				Value: float64(rtm.NumForcedGC),
			},
			{
				Type:  metrics.Gauge,
				Name:  "NumGC",
				Value: float64(rtm.NumGC),
			},
			{
				Type:  metrics.Gauge,
				Name:  "OtherSys",
				Value: float64(rtm.OtherSys),
			},
			{
				Type:  metrics.Gauge,
				Name:  "PauseTotalNs",
				Value: float64(rtm.PauseTotalNs),
			},
			{
				Type:  metrics.Gauge,
				Name:  "StackInuse",
				Value: float64(rtm.StackInuse),
			},
			{
				Type:  metrics.Gauge,
				Name:  "PauseTotalNs",
				Value: float64(rtm.PauseTotalNs),
			},
			{
				Type:  metrics.Gauge,
				Name:  "StackSys",
				Value: float64(rtm.StackSys),
			},
			{
				Type:  metrics.Gauge,
				Name:  "Sys",
				Value: float64(rtm.Sys),
			},
			{
				Type:  metrics.Gauge,
				Name:  "TotalAlloc",
				Value: float64(rtm.TotalAlloc),
			},
		}
		metricsCh <- metrics
	}()

	return metricsCh
}
