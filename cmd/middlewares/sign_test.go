package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SpaceSlow/execenv/cmd/config"
)

func TestWithSigning(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		key            string
		reqHeaderHash  string
		respHeaderHash string
		reqBody        []byte
		responseBody   []byte
		wantStatusCode int
	}{
		{
			name:           "post with non-empty body with correct hash header",
			method:         http.MethodPost,
			key:            "key",
			reqHeaderHash:  "84f1980ab8751660f844a7e6838b1324314b4aea330aa75ad7de0c62a1e807b3", // == hashSHA256("textkey")
			reqBody:        []byte("text"),
			respHeaderHash: "2c70e12b7a0646f92279f427c7b38e7334d8e5389cff167a1dc30e73f826b683", // == hashSHA256("key")
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "post with non-empty body with empty hash header",
			method:         http.MethodPost,
			key:            "key",
			reqBody:        []byte("text"),
			respHeaderHash: "2c70e12b7a0646f92279f427c7b38e7334d8e5389cff167a1dc30e73f826b683", // == hashSHA256("key")
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "post with non-empty body with incorrect hash header",
			method:         http.MethodPost,
			key:            "key",
			reqHeaderHash:  "incorrect-hash",
			reqBody:        []byte("text"),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "post with empty body with incorrect hash header",
			method:         http.MethodPost,
			key:            "key",
			reqHeaderHash:  "incorrect-hash",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "post with incorrect hash header (server config without key)",
			method:         http.MethodPost,
			reqHeaderHash:  "incorrect-hash",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "get with non-empty body with correct hash header",
			method:         http.MethodGet,
			key:            "key",
			reqHeaderHash:  "84f1980ab8751660f844a7e6838b1324314b4aea330aa75ad7de0c62a1e807b3", // == hashSHA256("textkey")
			reqBody:        []byte("text"),
			respHeaderHash: "2c70e12b7a0646f92279f427c7b38e7334d8e5389cff167a1dc30e73f826b683", // == hashSHA256("key")
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "get with non-empty body with empty hash header",
			method:         http.MethodGet,
			key:            "key",
			reqBody:        []byte("text"),
			respHeaderHash: "2c70e12b7a0646f92279f427c7b38e7334d8e5389cff167a1dc30e73f826b683", // == hashSHA256("key")
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "get with non-empty body with incorrect hash header",
			method:         http.MethodGet,
			key:            "key",
			reqHeaderHash:  "incorrect-hash",
			reqBody:        []byte("text"),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "get with empty body with incorrect hash header",
			method:         http.MethodPost,
			key:            "key",
			reqHeaderHash:  "incorrect-hash",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "get with incorrect hash header (server config without key)",
			method:         http.MethodGet,
			reqHeaderHash:  "incorrect-hash",
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{"test"}
			c, err := config.GetServerConfig()
			require.NoError(t, err)
			c.Key = tt.key

			handler := WithSigning(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write(tt.responseBody)
			}))

			var req *http.Request
			if tt.reqBody != nil {
				req, err = http.NewRequest(tt.method, "https://example.com", bytes.NewReader(tt.reqBody))
			} else {
				req, err = http.NewRequest(tt.method, "https://example.com", nil)
			}
			require.NoError(t, err)
			req.Header["Hash"] = []string{tt.reqHeaderHash}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
			if gotHeaderHash, ok := rr.Header()["Hash"]; tt.wantStatusCode == http.StatusOK && ok {
				require.Len(t, gotHeaderHash, 1)
				assert.Equal(t, tt.respHeaderHash, gotHeaderHash[0])
			}
		})
	}
}
