package storages

import "github.com/SpaceSlow/execenv/cmd/server/metrics"

type Storage interface {
	Add(metric *metrics.Metric) error
}
