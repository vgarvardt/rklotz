package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/vgarvardt/rklotz/pkg/server/rqctx"
)

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
		requestID := middleware.GetReqID(r.Context())
		requestLogger := m.logger.With(zap.String("request-id", requestID))

		newCtx := rqctx.SetID(r.Context(), requestID)
		newCtx = rqctx.SetLogger(newCtx, requestLogger)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
