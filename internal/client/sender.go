package client

import "github.com/SpaceSlow/execenv/internal/metrics"

type Sender interface {
	Send(metrics []metrics.Metric) error
}
