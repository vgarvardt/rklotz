package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

var static = []string{".css", ".js", ".png", ".jpg", ".jpeg", ".ico"}

// LoggerRequest implements chi/middleware.LogFormatter interface for requests logging
type LoggerRequest struct{}

// LoggerEntry implements chi/middleware.LogEntry interface for requests logging
type LoggerEntry struct {
	logger log.FieldLogger
	path   string
}

// NewLogEntry initiates the beginning of a new LogEntry per request.
func (l *LoggerRequest) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &LoggerEntry{path: r.URL.Path}

	entry.logger = log.WithFields(log.Fields{
		"method":      r.Method,
		"path":        r.URL.Path,
		"remote-addr": r.RemoteAddr,
		"user-agent":  r.UserAgent(),
	})

	entry.logger.Debug("Start serving request")

	return entry
}

// Write records the final log when a request completes
func (l *LoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.logger = l.logger.WithFields(log.Fields{
		"code":         status,
		"bytes_length": bytes,
		"elapsed_ms":   elapsed.String(),
	})

	msg := "Finished serving request"
	for i := range static {
		if strings.HasSuffix(l.path, static[i]) {
			l.logger.Debugln(msg)
			return
		}
	}

	l.logger.Infoln(msg)
}

// Panic records the final log when a request completes
func (l *LoggerEntry) Panic(v interface{}, stack []byte) {
	l.logger = l.logger.WithFields(log.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})

	l.logger.Errorln("Panic while serving request")
}
