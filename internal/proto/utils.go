package proto

import "github.com/SpaceSlow/execenv/internal/metrics"

func ConvertFromProto(m *Metric) (*metrics.Metric, error) {
	metric := &metrics.Metric{
		Name: m.Id,
	}

	switch m.MType {
	case MType_COUNTER:
		metric.Type = metrics.Counter
		metric.Value = m.Delta
	case MType_GAUGE:
		metric.Type = metrics.Gauge
		metric.Value = m.Value
	default:
		return nil, metrics.ErrIncorrectMetricTypeOrValue
	}
	return metric, nil
}

func ConvertToProto(m *metrics.Metric) (*Metric, error) {
	metric := &Metric{
		Id: m.Name,
	}

	switch m.Type {
	case metrics.Counter:
		metric.MType = MType_COUNTER
		delta, ok := m.Value.(int64)
		if !ok {
			return nil, metrics.ErrIncorrectMetricTypeOrValue
		}
		metric.Delta = delta
	case metrics.Gauge:
		metric.MType = MType_GAUGE
		value, ok := m.Value.(float64)
		if !ok {
			return nil, metrics.ErrIncorrectMetricTypeOrValue
		}
		metric.Value = value
	default:
		return nil, metrics.ErrIncorrectMetricTypeOrValue
	}
	return metric, nil
}
