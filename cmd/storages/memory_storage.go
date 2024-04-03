package storages

import (
	"sync"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

type MemStorage struct {
	mu      sync.Mutex
	metrics map[string]interface{}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: make(map[string]interface{})}
}

func (storage *MemStorage) Add(metric *metrics.Metric) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	switch metric.Type {
	case metrics.Counter:
		prevValue, _ := storage.metrics[metric.Name].(int64)
		value, _ := metric.Value.(int64)
		storage.metrics[metric.Name] = prevValue + value
	case metrics.Gauge:
		value, _ := metric.Value.(float64)
		storage.metrics[metric.Name] = value
	default:
		return &metrics.IncorrectMetricTypeOrValueError{}
	}
	return nil
}

func (storage *MemStorage) Get(name string) (metrics.Metric, bool) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	value, ok := storage.metrics[name]
	if !ok {
		return metrics.Metric{}, false
	}

	var metricType metrics.MetricType
	switch value.(type) {
	case float64:
		metricType = metrics.Gauge
	case int64:
		metricType = metrics.Counter
	default:
		return metrics.Metric{}, false
	}

	return metrics.Metric{
		Type:  metricType,
		Name:  name,
		Value: value,
	}, true
}
