package metrics

import "strconv"

type Metric struct {
	Name  string
	Type  MetricType
	Value string
}

func NewMetric(mT, name, value string) (*Metric, error) {
	metricType, err := parseMetricType(mT)
	if err != nil {
		return nil, &IncorrectMetricTypeOrValueError{}
	}

	switch metricType {
	case Counter:
		_, err = strconv.ParseInt(value, 10, 64)
	case Gauge:
		_, err = strconv.ParseFloat(value, 64)
	}

	if err != nil {
		return nil, &IncorrectMetricTypeOrValueError{}
	}

	return &Metric{name, metricType, value}, nil
}
