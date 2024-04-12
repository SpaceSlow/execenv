package metrics

type IncorrectMetricTypeOrValueError struct{}

func (e *IncorrectMetricTypeOrValueError) Error() string {
	return "incorrect metric type or value"
}

type EmptyMetricNameError struct{}

func (e *EmptyMetricNameError) Error() string {
	return "empty metric name"
}
