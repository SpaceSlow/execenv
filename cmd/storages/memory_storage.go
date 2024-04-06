package storages

import (
	"sync"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

type MemStorage struct {
	mu       sync.Mutex
	counters map[string]int64
	gauges   map[string]float64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{counters: make(map[string]int64), gauges: make(map[string]float64)}
}

func (storage *MemStorage) Add(metric *metrics.Metric) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	switch metric.Type {
	case metrics.Counter:
		prevValue := storage.counters[metric.Name]
		value, ok := metric.Value.(int64)
		if !ok {
			return &metrics.IncorrectMetricTypeOrValueError{}
		}
		storage.counters[metric.Name] = prevValue + value
	case metrics.Gauge:
		value, ok := metric.Value.(float64)
		if !ok {
			return &metrics.IncorrectMetricTypeOrValueError{}
		}
		storage.gauges[metric.Name] = value
	default:
		return &metrics.IncorrectMetricTypeOrValueError{}
	}
	return nil
}

func (storage *MemStorage) Get(metricType metrics.MetricType, name string) (metrics.Metric, bool) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	var value interface{}
	switch metricType {
	case metrics.Counter:
		value = storage.counters[name]
	case metrics.Gauge:
		value = storage.gauges[name]
	default:
		return metrics.Metric{}, false
	}

	return metrics.Metric{
		Type:  metricType,
		Name:  name,
		Value: value,
	}, true
}
