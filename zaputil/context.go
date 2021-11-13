package zaputil

import (
	"context"

	"go.uber.org/zap"
)

var contextKey = &struct{}{}

var defaultLogger = zap.NewNop()

// Extract pulls the logger from the provided context. If no logger is found, then the defaultLogger is returned.
func Extract(ctx context.Context) *zap.Logger {
	log := ctx.Value(contextKey)
	if log == nil {
		return defaultLogger
	}

	return log.(*zap.Logger)
}

// ToContext sets the logger on the provided context.
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}
