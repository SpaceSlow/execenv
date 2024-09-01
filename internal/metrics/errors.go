package metrics

import "errors"

var (
	ErrIncorrectMetricTypeOrValue = errors.New("incorrect metric type or value")
	ErrEmptyMetricName            = errors.New("empty metric name")
	ErrDecodePEMBlock             = errors.New("failed to decode PEM block containing public key")
)
