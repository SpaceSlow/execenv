package middlewares

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/SpaceSlow/execenv/internal/config"
)

var ErrHashEmptyBody = errors.New("hash empty body error")

func getHashBody(req *http.Request, key string) (string, error) {
	if req.Body == nil {
		return "", ErrHashEmptyBody
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	rb := bytes.NewReader(body)
	req.Body = io.NopCloser(rb)
	h := sha256.New()
	h.Write(append(body, []byte(key)...))

	return hex.EncodeToString(h.Sum(nil)), nil
}

// WithSigning middleware предназначенная для подписи данных и проверки подписи.
func WithSigning(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, err := config.GetServerConfig()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if headerHash := r.Header.Get("Hash"); headerHash != "none" && headerHash != "" && cfg.Key != "" {
			var hashSum string
			hashSum, err = getHashBody(r, cfg.Key)
			if err != nil && !errors.Is(err, ErrHashEmptyBody) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if hashSum != headerHash {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if cfg.Key == "" {
			next.ServeHTTP(w, r)
			return
		}

		l := httptest.NewRecorder()
		next.ServeHTTP(l, r)

		h := sha256.New()
		body, err := io.ReadAll(l.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.Write(append(body, []byte(cfg.Key)...))
		for key, values := range l.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.Header().Set("Hash", hex.EncodeToString(h.Sum(nil)))
		w.WriteHeader(l.Code)
		w.Write(body)
	})
}
