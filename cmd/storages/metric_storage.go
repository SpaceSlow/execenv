package storages

import "github.com/SpaceSlow/execenv/cmd/metrics"

type MetricStorage interface {
	Add(metric *metrics.Metric) (*metrics.Metric, error)
	Batch(metrics []metrics.Metric) ([]metrics.Metric, error)
	Get(metricType metrics.MetricType, name string) (*metrics.Metric, bool)
	List() []metrics.Metric
	Close() error
}
