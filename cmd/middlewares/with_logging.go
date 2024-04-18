package middlewares

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

var Log = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}

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

		Log.Info(
			"request/response",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", l.response.statusCode),
			zap.Int("size", l.response.size),
		)
	}

	return loggerHandleFunc
}
