package storages

import (
	"sync"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

type counters map[string]int64

func (c counters) Add(metric *metrics.Metric) (*metrics.Metric, error) {
	prevValue := c[metric.Name]
	value, ok := metric.Value.(int64)
	if !ok {
		return nil, metrics.ErrIncorrectMetricTypeOrValue
	}
	updValue := prevValue + value
	updMetric := metric.Copy()
	updMetric.Value = updValue
	c[metric.Name] = updValue

	return updMetric, nil
}

type gauges map[string]float64

func (g gauges) Add(metric *metrics.Metric) (*metrics.Metric, error) {
	value, ok := metric.Value.(float64)
	if !ok {
		return nil, metrics.ErrIncorrectMetricTypeOrValue
	}
	updMetric := metric.Copy()
	g[metric.Name] = value

	return updMetric, nil
}

type MemStorage struct {
	mu       sync.Mutex
	counters counters
	gauges   gauges
}

func NewMemStorage() *MemStorage {
	return &MemStorage{counters: make(map[string]int64), gauges: make(map[string]float64)}
}

func (storage *MemStorage) Add(metric *metrics.Metric) (*metrics.Metric, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	var (
		updMetric *metrics.Metric
		err       error
	)
	switch metric.Type {
	case metrics.Counter:
		updMetric, err = storage.counters.Add(metric)
	case metrics.Gauge:
		updMetric, err = storage.gauges.Add(metric)
	default:
		return nil, metrics.ErrIncorrectMetricTypeOrValue
	}
	if err != nil {
		return nil, err
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

func (storage *MemStorage) Batch(metricSlice []metrics.Metric) error {
	for _, metric := range metricSlice {
		_, err := storage.Add(&metric)
		if err != nil {
			return err
		}
	}
	return nil
}
