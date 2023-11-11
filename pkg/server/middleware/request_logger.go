package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"

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
		rqctx.GetLogger(r.Context()).Debug("Started request", slog.String("method", r.Method), slog.String("path", r.URL.Path))

		metrics := httpsnoop.CaptureMetrics(next, w, r)

		logEntry := rqctx.GetLogger(r.Context()).With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("host", r.Host),
			slog.String("request", r.RequestURI),
			slog.String("remote-addr", r.RemoteAddr),
			slog.String("referer", r.Referer()),
			slog.String("user-agent", r.UserAgent()),
			slog.Int("code", metrics.Code),
			slog.Int("duration-ms", int(metrics.Duration/time.Millisecond)),
			slog.String("duration-fmt", metrics.Duration.String()),
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
