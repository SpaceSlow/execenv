package middlewares

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
)

var KEY string

func getHashBody(req *http.Request, key string) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	rb := bytes.NewReader(body)
	req.Body = io.NopCloser(rb)
	h := sha256.New()
	h.Write(body)

	return hex.EncodeToString(h.Sum([]byte(key))), nil
}

func WithSigning(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if headerHash := r.Header.Get("Hash"); headerHash != "" && headerHash != "none" && KEY != "" {
			hashSum, err := getHashBody(r, KEY)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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

		h := sha256.New()
		body, err := io.ReadAll(l.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.Write(body)
		for key, values := range l.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.Header().Set("Hash", hex.EncodeToString(h.Sum([]byte(KEY))))
		w.WriteHeader(l.Code)
		w.Write(body)
	})
}
