package metrics

import (
	"strconv"
)

type Metric struct {
	Type  MetricType
	Name  string
	Value interface{}
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
