package metrics

import "fmt"

const (
	Counter = iota + 1
	Gauge
)

// MetricType хранит тип метрики Counter и Gauge.
type MetricType int

func (mt MetricType) isValid() bool {
	return mt >= Counter && mt <= Gauge
}

func (mt MetricType) String() string {
	metricTypes := []string{
		"counter",
		"gauge",
	}

	if !mt.isValid() {
		return fmt.Sprintf("MetricType(%d)", mt)
	}

	return metricTypes[mt-1]
}

// ParseMetricType возвращает тип метрики и ошибку, в случае невозможности разбора (поддерживает только lowercase).
func ParseMetricType(mType string) (MetricType, error) {
	switch mType {
	case "counter":
		return Counter, nil
	case "gauge":
		return Gauge, nil
	default:
		return MetricType(-1), ErrIncorrectMetricTypeOrValue
	}
}
