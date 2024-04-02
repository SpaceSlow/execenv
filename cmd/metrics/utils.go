package metrics

import "strings"

// ParseMetricFromPath `path` string represent "<metricName>/<value>"
func ParseMetricFromPath(path string, metricType MetricType) (*Metric, error) {
	metricFields := strings.FieldsFunc(path, func(r rune) bool { return r == '/' })

	if len(metricFields) == 0 {
		return nil, &EmptyMetricNameError{}
	}
	if len(metricFields) != 2 || !metricType.isValid() {
		return nil, &IncorrectMetricTypeOrValueError{}
	}

	return NewMetric(metricType, metricFields[0], metricFields[1])
}
