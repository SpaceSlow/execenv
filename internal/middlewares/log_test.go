package middlewares

import (
	"encoding/json"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SpaceSlow/execenv/internal/logger"
)

func TestWithLogging(t *testing.T) {
	expectedDuration := time.Second
	expectedMethod := []string{http.MethodGet, http.MethodPost}[rand.Intn(2)]

	file, err := os.CreateTemp(t.TempDir(), "temp.*.log")
	require.NoError(t, err)
	defer os.Remove(file.Name())
	defer file.Close()
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{file.Name()},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger.Log, err = cfg.Build()
	require.NoError(t, err)

	handler := WithLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(expectedDuration)
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
	}))

	req, err := http.NewRequest(expectedMethod, "https://example.org", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data, err := os.ReadFile(file.Name())
	require.NoError(t, err)
	require.NotEmpty(t, data)
	require.Contains(t, string(data), "uri")
	require.Contains(t, string(data), "method")
	require.Contains(t, string(data), "status")
	require.Contains(t, string(data), "size")
	require.Contains(t, string(data), "duration")

	logJSON := struct {
		URI      string        `json:"uri"`
		Method   string        `json:"method"`
		Status   int           `json:"status"`
		Size     int           `json:"size"`
		Duration time.Duration `json:"duration"`
	}{}

	err = json.Unmarshal(data, &logJSON)
	require.NoError(t, err)
	assert.Equal(t, expectedMethod, logJSON.Method)
	assert.Equal(t, 0, logJSON.Size)
	assert.Equal(t, http.StatusOK, logJSON.Status)
	assert.GreaterOrEqual(t, logJSON.Duration, expectedDuration)
}
