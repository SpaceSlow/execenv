package storages

import "github.com/SpaceSlow/execenv/cmd/metrics"

// MetricStorage является интерфейсом для хранения метрик.
type MetricStorage interface {
	Add(metric *metrics.Metric) (*metrics.Metric, error)
	Batch(metrics []metrics.Metric) error
	Get(metricType metrics.MetricType, name string) (*metrics.Metric, bool)
	List() []metrics.Metric
	Close() error
}
