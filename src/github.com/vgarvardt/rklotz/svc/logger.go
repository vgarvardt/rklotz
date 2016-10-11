package svc

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

var logger *log.Logger

func NewLiveLogger() *log.Logger {
	logger = &log.Logger{
		Out: os.Stderr,
		Formatter: &log.JSONFormatter{},
		Level: log.DebugLevel,
	}

	return logger
}
