package middleware

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pressly/chi/middleware"
)

type LoggerRequest struct{}

type LoggerEntry struct {
	logger log.FieldLogger
}

func (l *LoggerRequest) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &LoggerEntry{}

	entry.logger = log.WithFields(log.Fields{
		"method":      r.Method,
		"path":        r.URL.Path,
		"remote-addr": r.RemoteAddr,
		"user-agent":  r.UserAgent(),
	})

	entry.logger.Debug("Start serving request")

	return entry
}

func (l *LoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.logger = l.logger.WithFields(log.Fields{
		"code":         status,
		"bytes_length": bytes,
		"elapsed_ms":   elapsed.String(),
	})

	l.logger.Infoln("Finished serving request")
}

func (l *LoggerEntry) Panic(v interface{}, stack []byte) {
	l.logger = l.logger.WithFields(log.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}
