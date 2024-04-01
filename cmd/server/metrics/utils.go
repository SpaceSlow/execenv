package metrics

import "strings"

// ParseMetricFromURL `url` string represent "<metricName>/<value>"
func ParseMetricFromURL(url string, metricType MetricType) (*Metric, error) {
	metricFields := strings.FieldsFunc(url, func(r rune) bool { return r == '/' })

	if len(metricFields) == 0 {
		return nil, &EmptyMetricNameError{}
	}
	if len(metricFields) != 2 || !metricType.isValid() {
		return nil, &IncorrectMetricTypeOrValueError{}
	}

	return NewMetric(metricType, metricFields[0], metricFields[1])
}
