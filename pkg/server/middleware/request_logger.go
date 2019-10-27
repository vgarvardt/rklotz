package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"

	"github.com/vgarvardt/rklotz/pkg/server/rqctx"
)

var static = [...]string{".css", ".js", ".png", ".jpg", ".jpeg", ".ico"}

// RequestLogger is a middleware that logs request data
type RequestLogger struct{}

// NewRequestLogger creates new RequestLogger instance
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{}
}

// Handler is the request handler that logs request data
func (m *RequestLogger) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rqctx.GetLogger(r.Context()).Debug("Started request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

		metrics := httpsnoop.CaptureMetrics(next, w, r)

		logEntry := rqctx.GetLogger(r.Context()).With(
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("host", r.Host),
			zap.String("request", r.RequestURI),
			zap.String("remote-addr", r.RemoteAddr),
			zap.String("referer", r.Referer()),
			zap.String("user-agent", r.UserAgent()),
			zap.Int("code", metrics.Code),
			zap.Int("duration-ms", int(metrics.Duration/time.Millisecond)),
			zap.String("duration-fmt", metrics.Duration.String()),
		)

		if m.isStaticRequest(r) {
			logEntry.Debug("Finished serving request")
			return
		}

		logEntry.Info("Finished serving request")
	})
}

func (m *RequestLogger) isStaticRequest(r *http.Request) bool {
	for _, ext := range static {
		if strings.HasSuffix(r.URL.Path, ext) {
			return true
		}
	}

	return false
}
