package svc

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
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

func NewNullLogger() *log.Logger {
	logger, _ = test.NewNullLogger()
	return logger
}
