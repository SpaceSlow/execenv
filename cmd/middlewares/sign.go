package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

const HashHeaderField = "Hash"

var KEY string

func WithSigning(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if headerHash := r.Header.Get(HashHeaderField); headerHash != "" && headerHash != "none" && KEY != "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			hashSum := metrics.GetHash(body, KEY)
			if hashSum != headerHash {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if KEY == "" {
			next.ServeHTTP(w, r)
			return
		}

		l := httptest.NewRecorder()
		next.ServeHTTP(l, r)

		body, err := io.ReadAll(l.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		hash := metrics.GetHash(body, KEY)

		for key, values := range l.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.Header().Set(HashHeaderField, hash)
		w.WriteHeader(l.Code)
		w.Write(body)
	})
}
