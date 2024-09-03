package metrics

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func generateRandomMetric() Metric {
	mType := MetricType(rand.Intn(2) + 1)

	metric := Metric{
		Type: mType,
		Name: randStringBytes(rand.Intn(10) + 10),
	}

	switch mType {
	case Counter:
		metric.Value = rand.Int63n(1000)
	case Gauge:
		metric.Value = rand.Float64()
	}
	return metric
}

func TestFanIn(t *testing.T) {
	metricsChs := make([]chan []Metric, 10)
	for i := range metricsChs {
		metricsChs[i] = make(chan []Metric, 1)
	}

	expected := make([]Metric, 100)
	for i := range expected {
		expected[i] = generateRandomMetric()
	}

	for i := 0; i < len(metricsChs); i++ {
		metricsChs[i] <- expected[i*10 : i*10+10]
		close(metricsChs[i])
	}

	got := make([]Metric, 0)
	for slice := range FanIn(metricsChs...) {
		got = append(got, slice...)
	}

	assert.ElementsMatch(t, expected, got)
}
