package metrics

import "fmt"

const (
	Counter = iota + 1
	Gauge
)

type MetricType int

func (mt MetricType) isValid() bool {
	return mt >= Counter && mt <= Gauge
}

func (mt MetricType) String() string {
	metricTypes := [...]string{"counter", "gauge"}
	if !mt.isValid() {
		return fmt.Sprintf("MetricType(%d)", mt)
	}
	return metricTypes[mt-1]
}

func parseMetricType(s string) (MetricType, error) {
	switch s {
	case "counter":
		return Counter, nil
	case "gauge":
		return Gauge, nil
	default:
		return -1, &UnknownMetricTypeError{}
	}
}
