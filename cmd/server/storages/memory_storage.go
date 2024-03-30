package storages

import (
	"strconv"

	"github.com/SpaceSlow/execenv/cmd/server/metrics"
)

type MemStorage struct {
	metrics map[string]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: make(map[string]string)}
}

func (storage *MemStorage) Add(metric *metrics.Metric) error {
	switch metric.Type {
	case metrics.Counter:
		{
			prevValue, _ := strconv.ParseInt(storage.metrics[metric.Name], 10, 64)
			value, _ := strconv.ParseInt(metric.Value, 10, 64)
			storage.metrics[metric.Name] = strconv.FormatInt(prevValue+value, 10)

		}
	case metrics.Gauge:
		storage.metrics[metric.Name] = metric.Value
	default:
		return &metrics.UnknownMetricTypeError{}
	}
	return nil
}
