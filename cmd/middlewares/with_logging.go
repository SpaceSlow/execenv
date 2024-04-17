package middlewares

import (
	"log/slog"
	"net/http"
	"time"
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

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	loggerHandleFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := loggingResponseWriter{
			ResponseWriter: w,
			response: &Response{
				statusCode: 0,
				size:       0,
			},
		}
		h(&l, r)

		duration := time.Since(start)

		slog.Default().Info(
			"request/response",
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", l.response.statusCode,
			"size", l.response.size,
		)
	}

	return loggerHandleFunc
}
