// Copyright (C) 2021 Mya Pitzeruse
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

// HashicorpStdLogger
// Deprecated.
func HashicorpStdLogger(logger *zap.Logger) *log.Logger {
	return HashiCorpStdLogger(logger)
}

// HashiCorpStdLogger wraps the provided logger with a golang logger to log messages at the appropriate level using the
// HashiCorp log format. Useful for replacing serf and membership loggers.
func HashiCorpStdLogger(logger *zap.Logger) *log.Logger {
	w := &hashiWriter{
		logger: logger.WithOptions(zap.AddCallerSkip(3)),
	}

	return log.New(w, "", 0)
}
