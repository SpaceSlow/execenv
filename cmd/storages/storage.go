package storages

import (
	"github.com/SpaceSlow/execenv/cmd/metrics"
)

type Storage interface {
	Add(metric *metrics.Metric) error
	Get(name string) (metrics.Metric, bool)
}
