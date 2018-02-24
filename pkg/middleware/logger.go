package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// RequestLoggerKey is the key that holds th unique request logger in a request context.
const RequestLoggerKey = contextKey("RequestLogger")

// Logger is a middleware that injects logger with request ID into the context of each request.
type Logger struct {
	logger *zap.Logger
}

// NewLogger creates new Logger instance
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}
}

// Handler is the request handler that creates logger instance for each request with corresponding request ID.
func (m *Logger) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestLogger := m.logger.With(zap.String("request-id", middleware.GetReqID(r.Context())))
		ctx := context.WithValue(r.Context(), RequestLoggerKey, requestLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestLogger returns a request logger from the given context if one is present.
func GetRequestLogger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		panic("Can not get request logger from empty context")
	}
	if requestLogger, ok := ctx.Value(RequestLoggerKey).(*zap.Logger); ok {
		return requestLogger
	}

	return nil
}
