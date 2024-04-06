package metrics

import (
	"fmt"
	"strconv"
)

type Metric struct {
	Type  MetricType
	Name  string
	Value interface{}
}

func (m *Metric) String() string {
	return fmt.Sprintf("%s = %v (%s)", m.Name, m.Value, m.Type)
}

func (m *Metric) ValueAsString() string {
	switch m.Type {
	case Counter:
		return strconv.FormatInt(m.Value.(int64), 10)
	case Gauge:
		return strconv.FormatFloat(m.Value.(float64), 'f', -1, 64)
	default:
		return ""
	}
}

func NewMetric(metricType MetricType, name, value string) (*Metric, error) {
	var err error
	var val interface{}
	switch metricType {
	case Counter:
		val, err = strconv.ParseInt(value, 10, 64)
	case Gauge:
		val, err = strconv.ParseFloat(value, 64)
	}

	if err != nil {
		return nil, &IncorrectMetricTypeOrValueError{}
	}

	return &Metric{metricType, name, val}, nil
}
