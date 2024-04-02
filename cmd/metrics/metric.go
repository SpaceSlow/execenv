package metrics

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
	"strconv"
)

const Url = "http://localhost:8080/update/"
const Slash = "/"

var PollCount int64

type Metric struct {
	Type  MetricType
	Name  string
	Value interface{}
}

func NewMetric(metricType MetricType, name, value string) (*Metric, error) {
	var err error
	var val interface{}
	switch metricType {
	case Counter:
		val, err = strconv.ParseInt(value, 10, 64)
	case Gauge:
		val, err = strconv.ParseFloat(value, 64)
	}

	if err != nil {
		return nil, &IncorrectMetricTypeOrValueError{}
	}

	return &Metric{metricType, name, val}, nil
}

func SendMetrics(metrics []Metric) {
	for _, metric := range metrics {
		_ = SendMetric(metric)
		PollCount = 0
	}
}

func SendMetric(metric Metric) error {
	var err error
	_, err = http.Post(
		fmt.Sprint(Url, metric.Type.String(), Slash, metric.Name, Slash, metric.Value),
		"text/plain",
		nil)
	return err
}

func GetMetrics() []Metric {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	PollCount++

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
		{
			Type:  Gauge,
			Name:  "RandomValue",
			Value: rand.Float64(),
		},
		{
			Type:  Counter,
			Name:  "PollCount",
			Value: PollCount,
		},
	}

	return metrics
}
