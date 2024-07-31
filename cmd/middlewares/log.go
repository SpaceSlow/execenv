package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/SpaceSlow/execenv/cmd/logger"
)

var _ http.ResponseWriter = (*loggingResponseWriter)(nil)

type response struct {
	statusCode int
	size       int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	response *response
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

// WithLogging middleware предназначенная для логирования запросов пользователей.
// В логи попадает следующая информация: uri, метод запроса, продолжительность обработки, статус ответа и размер ответа.
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := loggingResponseWriter{
			ResponseWriter: w,
			response: &response{
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
