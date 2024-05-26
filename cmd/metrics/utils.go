package metrics

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func ParseMetricType(mType string) (MetricType, error) {
	switch mType {
	case "counter":
		return Counter, nil
	case "gauge":
		return Gauge, nil
	default:
		return MetricType(-1), ErrIncorrectMetricTypeOrValue
	}
}

func RetryFunc(f func() error, delays []time.Duration) chan error {
	errorCh := make(chan error)

	go func() {
		defer close(errorCh)
		var err error
		for attempt := 0; attempt < len(delays); attempt++ {
			if err = f(); err == nil {
				errorCh <- nil
				return
			}
			<-time.After(delays[attempt])
		}
		errorCh <- err
	}()

	return errorCh
}

func newCompressedRequest(method, url string, data []byte) (*http.Request, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	req, err := http.NewRequest(method, url, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func GetMetrics(metricsCh chan []Metric) {
	metricSlice := GetRuntimeMetrics()
	metricSlice = append(
		metricSlice,
		Metric{
			Type:  Gauge,
			Name:  "RandomValue",
			Value: rand.Float64(),
		},
	)
	metricsCh <- metricSlice
}

func GetRuntimeMetrics() []Metric {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	metrics := []Metric{
		{
			Type:  Gauge,
			Name:  "Alloc",
			Value: float64(rtm.Alloc),
		},
		{
			Type:  Gauge,
			Name:  "BuckHashSys",
			Value: float64(rtm.BuckHashSys),
		},
		{
			Type:  Gauge,
			Name:  "Frees",
			Value: float64(rtm.Frees),
		},
		{
			Type:  Gauge,
			Name:  "GCCPUFraction",
			Value: float64(rtm.GCCPUFraction),
		},
		{
			Type:  Gauge,
			Name:  "GCSys",
			Value: float64(rtm.GCSys),
		},
		{
			Type:  Gauge,
			Name:  "HeapAlloc",
			Value: float64(rtm.HeapAlloc),
		},
		{
			Type:  Gauge,
			Name:  "HeapIdle",
			Value: float64(rtm.HeapIdle),
		},
		{
			Type:  Gauge,
			Name:  "HeapInuse",
			Value: float64(rtm.HeapInuse),
		},
		{
			Type:  Gauge,
			Name:  "HeapObjects",
			Value: float64(rtm.HeapObjects),
		},
		{
			Type:  Gauge,
			Name:  "HeapReleased",
			Value: float64(rtm.HeapReleased),
		},
		{
			Type:  Gauge,
			Name:  "HeapSys",
			Value: float64(rtm.HeapSys),
		},
		{
			Type:  Gauge,
			Name:  "LastGC",
			Value: float64(rtm.LastGC),
		},
		{
			Type:  Gauge,
			Name:  "Lookups",
			Value: float64(rtm.Lookups),
		},
		{
			Type:  Gauge,
			Name:  "MCacheInuse",
			Value: float64(rtm.MCacheInuse),
		},
		{
			Type:  Gauge,
			Name:  "MCacheSys",
			Value: float64(rtm.MCacheSys),
		},
		{
			Type:  Gauge,
			Name:  "MSpanInuse",
			Value: float64(rtm.MSpanInuse),
		},
		{
			Type:  Gauge,
			Name:  "MSpanSys",
			Value: float64(rtm.MSpanSys),
		},
		{
			Type:  Gauge,
			Name:  "Mallocs",
			Value: float64(rtm.Mallocs),
		},
		{
			Type:  Gauge,
			Name:  "NextGC",
			Value: float64(rtm.NextGC),
		},
		{
			Type:  Gauge,
			Name:  "NumForcedGC",
			Value: float64(rtm.NumForcedGC),
		},
		{
			Type:  Gauge,
			Name:  "NumGC",
			Value: float64(rtm.NumGC),
		},
		{
			Type:  Gauge,
			Name:  "OtherSys",
			Value: float64(rtm.OtherSys),
		},
		{
			Type:  Gauge,
			Name:  "PauseTotalNs",
			Value: float64(rtm.PauseTotalNs),
		},
		{
			Type:  Gauge,
			Name:  "StackInuse",
			Value: float64(rtm.StackInuse),
		},
		{
			Type:  Gauge,
			Name:  "PauseTotalNs",
			Value: float64(rtm.PauseTotalNs),
		},
		{
			Type:  Gauge,
			Name:  "StackSys",
			Value: float64(rtm.StackSys),
		},
		{
			Type:  Gauge,
			Name:  "Sys",
			Value: float64(rtm.Sys),
		},
		{
			Type:  Gauge,
			Name:  "TotalAlloc",
			Value: float64(rtm.TotalAlloc),
		},
	}

	return metrics
}
