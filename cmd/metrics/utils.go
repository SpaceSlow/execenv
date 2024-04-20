package metrics

import (
	"bytes"
	"net/http"
	"runtime"
)

func ParseMetricType(mType string) (MetricType, error) {
	switch mType {
	case "counter":
		return Counter, nil
	case "gauge":
		return Gauge, nil
	default:
		return MetricType(-1), &IncorrectMetricTypeOrValueError{}
	}
}

func SendMetrics(url string, metrics []Metric) {
	for _, metric := range metrics {
		_ = SendMetric(url, metric)
	}
}

func SendMetric(url string, metric Metric) error {
	jsonMetric, err := metric.MarshalJSON()
	if err != nil {
		return err
	}
	readerMetric := bytes.NewReader(jsonMetric)

	res, err := http.Post(
		url,
		"application/json",
		readerMetric,
	)
	if err != nil {
		return err
	}
	return res.Body.Close()
}

func GetRuntimeMetrics() []Metric {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	metrics := []Metric{
		{
			Type:  Gauge,
			Name:  "Alloc",
			Value: rtm.Alloc,
		},
		{
			Type:  Gauge,
			Name:  "BuckHashSys",
			Value: rtm.BuckHashSys,
		},
		{
			Type:  Gauge,
			Name:  "Frees",
			Value: rtm.Frees,
		},
		{
			Type:  Gauge,
			Name:  "GCCPUFraction",
			Value: rtm.GCCPUFraction,
		},
		{
			Type:  Gauge,
			Name:  "GCSys",
			Value: rtm.GCSys,
		},
		{
			Type:  Gauge,
			Name:  "HeapAlloc",
			Value: rtm.HeapAlloc,
		},
		{
			Type:  Gauge,
			Name:  "HeapIdle",
			Value: rtm.HeapIdle,
		},
		{
			Type:  Gauge,
			Name:  "HeapInuse",
			Value: rtm.HeapInuse,
		},
		{
			Type:  Gauge,
			Name:  "HeapObjects",
			Value: rtm.HeapObjects,
		},
		{
			Type:  Gauge,
			Name:  "HeapReleased",
			Value: rtm.HeapReleased,
		},
		{
			Type:  Gauge,
			Name:  "HeapSys",
			Value: rtm.HeapSys,
		},
		{
			Type:  Gauge,
			Name:  "LastGC",
			Value: rtm.LastGC,
		},
		{
			Type:  Gauge,
			Name:  "Lookups",
			Value: rtm.Lookups,
		},
		{
			Type:  Gauge,
			Name:  "MCacheInuse",
			Value: rtm.MCacheInuse,
		},
		{
			Type:  Gauge,
			Name:  "MCacheSys",
			Value: rtm.MCacheSys,
		},
		{
			Type:  Gauge,
			Name:  "MSpanInuse",
			Value: rtm.MSpanInuse,
		},
		{
			Type:  Gauge,
			Name:  "MSpanSys",
			Value: rtm.MSpanSys,
		},
		{
			Type:  Gauge,
			Name:  "Mallocs",
			Value: rtm.Mallocs,
		},
		{
			Type:  Gauge,
			Name:  "NextGC",
			Value: rtm.NextGC,
		},
		{
			Type:  Gauge,
			Name:  "NumForcedGC",
			Value: rtm.NumForcedGC,
		},
		{
			Type:  Gauge,
			Name:  "NumGC",
			Value: rtm.NumGC,
		},
		{
			Type:  Gauge,
			Name:  "OtherSys",
			Value: rtm.OtherSys,
		},
		{
			Type:  Gauge,
			Name:  "PauseTotalNs",
			Value: rtm.PauseTotalNs,
		},
		{
			Type:  Gauge,
			Name:  "StackInuse",
			Value: rtm.StackInuse,
		},
		{
			Type:  Gauge,
			Name:  "PauseTotalNs",
			Value: rtm.PauseTotalNs,
		},
		{
			Type:  Gauge,
			Name:  "StackSys",
			Value: rtm.StackSys,
		},
		{
			Type:  Gauge,
			Name:  "Sys",
			Value: rtm.Sys,
		},
		{
			Type:  Gauge,
			Name:  "TotalAlloc",
			Value: rtm.TotalAlloc,
		},
	}

	return metrics
}
