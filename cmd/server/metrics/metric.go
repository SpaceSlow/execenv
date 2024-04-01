package metrics

import "strconv"

type Metric struct {
	Type  MetricType
	Name  string
	Value string
}

func NewMetric(metricType MetricType, name, value string) (*Metric, error) {
	var err error
	switch metricType {
	case Counter:
		_, err = strconv.ParseInt(value, 10, 64)
	case Gauge:
		_, err = strconv.ParseFloat(value, 64)
	}

	if err != nil {
		return nil, &IncorrectMetricTypeOrValueError{}
	}

	return &Metric{metricType, name, value}, nil
}
