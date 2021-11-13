package zaputil

import (
	"io"
	"log"
	"strings"

	"go.uber.org/zap"
)

// hashiWriter parses lines written from the Hashicorp ecosystem to proper log levels for zap. The format is as follows:
// [LEVEL] system: message.
type hashiWriter struct {
	logger *zap.Logger
}

func (w *hashiWriter) Write(p []byte) (n int, err error) {
	line := string(p)

	parts := strings.SplitN(line, " ", 3)
	for len(parts) < 3 {
		parts = append(parts, "")
	}

	level := strings.TrimSuffix(strings.TrimPrefix(parts[0], "["), "]")
	name := strings.TrimSuffix(parts[1], ":")
	msg := strings.TrimSpace(parts[2])

	fields := []zap.Field{
		zap.String("name", name),
	}

	switch level {
	case "DEBUG", "debug":
		w.logger.Debug(msg, fields...)
	case "INFO", "info":
		w.logger.Info(msg, fields...)
	case "WARN", "warn":
		w.logger.Warn(msg, fields...)
	case "ERROR", "error":
		w.logger.Error(msg, fields...)
	}

	return len(p), nil
}

var _ io.Writer = &hashiWriter{}

// HashicorpStdLogger wraps the provided logger with a golang logger to log messages at the appropriate level using the
// Hashicorp log format. Useful for replacing serf and membership loggers.
func HashicorpStdLogger(logger *zap.Logger) *log.Logger {
	w := &hashiWriter{
		logger: logger.WithOptions(zap.AddCallerSkip(3)),
	}

	return log.New(w, "", 0)
}
