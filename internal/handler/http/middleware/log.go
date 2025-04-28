package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

type (
	// responseWriterWrapper обёртка
	responseWriterWrapper struct {
		http.ResponseWriter
		statusCode int
	}
)

// WriteHeader перехватывает вызов WriteHeader
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// LogHandlerMiddleware логирование
func LogHandlerMiddleware(l *zap.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

			handler.ServeHTTP(wrappedWriter, r)
			l.Info("method log", zap.String("method", r.Method),
				zap.String("Request", r.URL.Path), zap.Int("status", wrappedWriter.statusCode))
		})
	}
}
