package metrics

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// ParseMetricFromPath `path` string represent "<metricName>/<value>"
func ParseMetricFromPath(path string, metricType MetricType) (*Metric, error) {
	metricFields := strings.FieldsFunc(path, func(r rune) bool { return r == '/' })

	if len(metricFields) == 0 {
		return nil, &EmptyMetricNameError{}
	}
	if len(metricFields) != 2 || !metricType.isValid() {
		return nil, &IncorrectMetricTypeOrValueError{}
	}

	return NewMetric(metricType, metricFields[0], metricFields[1])
}

func SendMetrics(url string, metrics []Metric) {
	for _, metric := range metrics {
		_ = SendMetric(url, metric)
	}
}

func SendMetric(url string, metric Metric) error {
	var err error
	_, err = http.Post(
		fmt.Sprintf(`%s/%s/%s/%v`, url, metric.Type.String(), metric.Name, metric.Value),
		"text/plain",
		nil)
	return err
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
