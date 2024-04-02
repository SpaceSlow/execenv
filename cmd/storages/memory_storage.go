package storages

import (
	"strconv"
	"sync"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

type MemStorage struct {
	mu      sync.Mutex
	metrics map[string]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: make(map[string]string)}
}

func (storage *MemStorage) Add(metric *metrics.Metric) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	switch metric.Type {
	case metrics.Counter:
		prevValue, _ := strconv.ParseInt(storage.metrics[metric.Name], 10, 64)
		value, _ := metric.Value.(int64)
		storage.metrics[metric.Name] = strconv.FormatInt(prevValue+value, 10)
	case metrics.Gauge:
		storage.metrics[metric.Name], _ = metric.Value.(string)
	default:
		return &metrics.IncorrectMetricTypeOrValueError{}
	}
	return nil
}
