package worker

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricWorkers_Err(t *testing.T) {
	cleanArgs()
	mw, err := NewMetricWorkers()
	require.NoError(t, err)
	mw.errorsCh <- errors.New("some error")
	mw.Close()
	assert.Error(t, <-mw.Err())
	assert.Nil(t, <-mw.Err()) // checking for the absence of a deadlock
}

func TestMetricWorkers_getGopsutilMetrics(t *testing.T) {
	cleanArgs()
	mw, err := NewMetricWorkers()
	require.NoError(t, err)
	metrics := <-mw.getGopsutilMetrics()
	assert.Greater(t, len(metrics), 0)

	for _, m := range metrics {
		assert.NotNil(t, m)
	}
}

func TestMetricWorkers_getRuntimeMetrics(t *testing.T) {
	cleanArgs()
	mw, err := NewMetricWorkers()
	require.NoError(t, err)
	metrics := <-mw.getRuntimeMetrics()
	assert.Greater(t, len(metrics), 0)

	for _, m := range metrics {
		assert.NotNil(t, m)
	}
}

func cleanArgs() {
	os.Args = []string{"program"}
}
