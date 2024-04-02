package metrics

const (
	Counter = iota + 1
	Gauge
)

type MetricType int

func (mt MetricType) isValid() bool {
	return mt >= Counter && mt <= Gauge
}
