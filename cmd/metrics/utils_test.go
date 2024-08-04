package metrics

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newCompressedRequest(t *testing.T) {
	correctReq, err := newCompressedRequest(http.MethodPost, "/api/updates", []byte(`[{"id":"PollCount", "delta": 5, "type":"counter"}]`))
	require.NoError(t, err)
	assert.Equal(t, correctReq.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, correctReq.Header.Get("Content-Encoding"), "gzip")
}
