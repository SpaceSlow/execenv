package middlewares

import (
	"bytes"
	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 ,.:"

func TestWithCompressing_compressedResponse(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		wantStatusCode int
	}{
		{
			name:           "post with medium request body text",
			request:        generateRequestWithLargeText(300),
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "post with nil request body",
			request:        requestWithNilBody(),
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "post with large request body text",
			request:        generateRequestWithLargeText(1500),
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := WithCompressing(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Body == nil {
					return
				}
				data, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "cannot read request body", http.StatusInternalServerError)
				}
				w.Header().Add("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			}))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, tt.request)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
			if tt.request.Body != nil {
				assert.Equal(t, CompressionAlgorithm, rr.Header().Get("Content-Encoding"))
				assert.Less(t, int64(rr.Body.Len()), tt.request.ContentLength)
			}
		})
	}
}

func TestWithCompressing_compressedRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		wantStatusCode int
	}{
		{
			name:           "post with compressed request body text",
			request:        generateCompressRequest(),
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := WithCompressing(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Body == nil {
					return
				}
				data, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "cannot read request body", http.StatusInternalServerError)
				}
				w.Write(data)
			}))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, tt.request)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
			if tt.request.Body != nil {
				assert.Less(t, tt.request.ContentLength, int64(rr.Body.Len()))
			}
		})
	}
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func generateCompressRequest() *http.Request {
	data, _ := metrics.Compress([]byte(randStringBytes(rand.Intn(500) + 500)))
	req, _ := http.NewRequest([]string{http.MethodGet, http.MethodPost}[rand.Intn(2)], "https://example.com", bytes.NewReader(data))
	req.Header.Add("Content-Encoding", CompressionAlgorithm)
	return req
}

func generateRequestWithLargeText(length int) *http.Request {
	req, _ := http.NewRequest(
		[]string{http.MethodGet, http.MethodPost}[rand.Intn(2)],
		"https://example.com",
		bytes.NewReader([]byte(randStringBytes(length))),
	)
	req.Header.Add("Accept-Encoding", CompressionAlgorithm)
	return req
}

func requestWithNilBody() *http.Request {
	req, _ := http.NewRequest([]string{http.MethodGet, http.MethodPost}[rand.Intn(2)], "https://example.com", nil)
	req.Header.Add("Accept-Encoding", CompressionAlgorithm)
	return req
}
