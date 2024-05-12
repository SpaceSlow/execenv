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

func (storage *MemStorage) Add(metric *metrics.Metric) (*metrics.Metric, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	var updMetric *metrics.Metric
	switch metric.Type {
	case metrics.Counter:
		prevValue := storage.counters[metric.Name]
		value, ok := metric.Value.(int64)
		if !ok {
			return nil, metrics.ErrIncorrectMetricTypeOrValue
		}
		updValue := prevValue + value
		updMetric = metric.Copy()
		updMetric.Value = updValue
		storage.counters[metric.Name] = updValue
	case metrics.Gauge:
		value, ok := metric.Value.(float64)
		if !ok {
			return nil, metrics.ErrIncorrectMetricTypeOrValue
		}
		updMetric = metric.Copy()
		storage.gauges[metric.Name] = value
	default:
		return nil, metrics.ErrIncorrectMetricTypeOrValue
	}
	return updMetric, nil
}

func (storage *MemStorage) Get(metricType metrics.MetricType, name string) (*metrics.Metric, bool) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	var value interface{}
	var ok bool
	switch metricType {
	case metrics.Counter:
		value, ok = storage.counters[name]
	case metrics.Gauge:
		value, ok = storage.gauges[name]
	default:
		return nil, false
	}
	if !ok {
		return nil, ok
	}

	return &metrics.Metric{
		Type:  metricType,
		Name:  name,
		Value: value,
	}, true
}

func (storage *MemStorage) List() []metrics.Metric {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	metricSlice := make([]metrics.Metric, 0, len(storage.counters)+len(storage.gauges))

	for name, value := range storage.counters {
		metricSlice = append(metricSlice, metrics.Metric{
			Type:  metrics.Counter,
			Name:  name,
			Value: value,
		})
	}
	for name, value := range storage.gauges {
		metricSlice = append(metricSlice, metrics.Metric{
			Type:  metrics.Gauge,
			Name:  name,
			Value: value,
		})
	}

	return metricSlice
}

func (storage *MemStorage) Close() error {
	return nil
}

func (storage *MemStorage) Batch(metricSlice []metrics.Metric) ([]metrics.Metric, error) {
	updMetrics := make([]metrics.Metric, 0, len(metricSlice))
	for _, metric := range metricSlice {
		updMetric, err := storage.Add(&metric)
		if err != nil {
			return updMetrics, err
		}
		updMetrics = append(updMetrics, *updMetric)
	}
	return updMetrics, nil
}
