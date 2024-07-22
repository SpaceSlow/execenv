package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/SpaceSlow/execenv/cmd/logger"
)

type Response struct {
	statusCode int
	size       int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	response *Response
}

func (l loggingResponseWriter) Write(bytes []byte) (int, error) {
	size, err := l.ResponseWriter.Write(bytes)
	l.response.size += size

	return size, err
}

func (l loggingResponseWriter) WriteHeader(statusCode int) {
	l.response.statusCode = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := loggingResponseWriter{
			ResponseWriter: w,
			response: &Response{
				statusCode: 0,
				size:       0,
			},
		}
		next.ServeHTTP(&l, r)

		duration := time.Since(start)

		logger.Log.Info(
			"request/response",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", l.response.statusCode),
			zap.Int("size", l.response.size),
		)
	})
}
