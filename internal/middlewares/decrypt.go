package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"

	"github.com/SpaceSlow/execenv/internal/config"
)

// WithDecryption middleware предназначенная для расшифрования данных с агента.
func WithDecryption(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, err := config.GetServerConfig()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if cfg.PrivateKey() == nil || r.Body == nil {
			next.ServeHTTP(w, r)
			return
		}

		encryptedData, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, cfg.PrivateKey(), encryptedData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rb := bytes.NewReader(decryptedData)
		r.Body = io.NopCloser(rb)

		next.ServeHTTP(w, r)
	})
}
