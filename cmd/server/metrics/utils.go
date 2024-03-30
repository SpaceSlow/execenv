package metrics

import "strings"

func ParseMetricFromURL(url string) (*Metric, error) {
	metricFields := strings.FieldsFunc(url, func(r rune) bool {
		return r == '/'
	})

	if len(metricFields) == 0 {
		return nil, &IncorrectMetricTypeOrValueError{}
	}
	metricType, err := parseMetricType(metricFields[0])

	if err == nil && len(metricFields) == 1 {
		return nil, &EmptyMetricNameError{}
	} else if len(metricFields) != 3 || !metricType.isValid() {
		return nil, &IncorrectMetricTypeOrValueError{}
	}

	return NewMetric(metricFields[0], metricFields[1], metricFields[2])
}
