package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/SpaceSlow/execenv/cmd/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func writeRandomPrivateKeyToFile(filename string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	pemPublicKey := pem.EncodeToMemory(block)

	return os.WriteFile(filename, pemPublicKey, 0600)
}

func TestWithDecryption(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		reqBody        []byte
		wantBody       []byte
		wantStatusCode int
	}{
		{
			name:           "post with non-empty request body",
			method:         http.MethodPost,
			reqBody:        []byte("text"),
			wantBody:       []byte("text"),
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "post with nil request body",
			method:         http.MethodPost,
			wantBody:       []byte(""),
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "post with empty request body",
			method:         http.MethodPost,
			reqBody:        make([]byte, 0),
			wantBody:       []byte(""),
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp(t.TempDir(), "private.*.key")
			require.NoError(t, err)
			defer os.Remove(file.Name())
			err = writeRandomPrivateKeyToFile(file.Name())
			require.NoError(t, err)

			os.Args = []string{"test", ""}
			os.Setenv("CRYPTO_KEY", file.Name())

			cfg, err := config.GetServerConfig()
			require.NoError(t, err)
			require.NotNil(t, cfg.PrivateKey())

			handler := WithDecryption(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Body == nil {
					return
				}
				data, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "cannot read request body", http.StatusInternalServerError)
				}
				w.Write(data)
			}))

			var req *http.Request
			if tt.reqBody != nil {
				encryptData, err := rsa.EncryptPKCS1v15(rand.Reader, &cfg.PrivateKey().PublicKey, tt.reqBody)
				require.NoError(t, err)
				req, err = http.NewRequest(tt.method, "https://example.com", bytes.NewReader(encryptData))
			} else {
				req, err = http.NewRequest(tt.method, "https://example.com", nil)
			}
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
			responseBody, err := io.ReadAll(rr.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody, responseBody)
		})
	}
}
