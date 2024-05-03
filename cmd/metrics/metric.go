package metrics

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Metric struct {
	Type  MetricType
	Name  string
	Value interface{}
}

type JSONMetric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metric) MarshalJSON() ([]byte, error) {
	metric := JSONMetric{
		ID:    m.Name,
		MType: m.Type.String(),
	}

	switch m.Type {
	case Counter:
		delta, ok := m.Value.(int64)
		if !ok {
			return nil, &IncorrectMetricTypeOrValueError{}
		}
		metric.Delta = &delta
	case Gauge:
		value, ok := m.Value.(float64)
		if !ok {
			return nil, &IncorrectMetricTypeOrValueError{}
		}
		metric.Value = &value
	default:
		return nil, &IncorrectMetricTypeOrValueError{}
	}

	return json.Marshal(metric)
}

func (m *Metric) UnmarshalJSON(data []byte) error {
	var metric JSONMetric
	if json.Unmarshal(data, &metric) != nil {
		return &IncorrectMetricTypeOrValueError{}
	}

	var mType MetricType
	if t, err := ParseMetricType(metric.MType); err != nil {
		return err
	} else {
		mType = t
	}
	m.Name = metric.ID
	m.Type = mType
	switch mType {
	case Counter:
		if metric.Delta != nil {
			m.Value = *metric.Delta
		}
	case Gauge:
		if metric.Value != nil {
			m.Value = *metric.Value
		}
	}
	if m.Value == nil {
		return &IncorrectMetricTypeOrValueError{}
	}
	return nil
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
