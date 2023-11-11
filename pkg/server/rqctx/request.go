package rqctx

import (
	"context"
	"log/slog"
)

type requestCtxKey int

const (
	requestIDKey requestCtxKey = iota
	requestLoggerKey
)

// SetID sets request ID to the request context
func SetID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// SetLogger sets logger to the request context
func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, requestLoggerKey, logger)
}

// GetLogger returns a request logger from the given context if one is present
func GetLogger(ctx context.Context) *slog.Logger {
	if ctx == nil {
		panic("Can not get request logger from empty context")
	}
	if requestLogger, ok := ctx.Value(requestLoggerKey).(*slog.Logger); ok {
		return requestLogger
	}

	return nil
}

// GetID returns a request ID from the given context if one is present
func GetID(ctx context.Context) string {
	if ctx == nil {
		panic("Can not get request ID from empty context")
	}

	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}

	return ""
}
