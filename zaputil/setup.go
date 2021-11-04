package zaputil

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DefaultConfig returns the default configuration for zap to use. By default, it logs at an info level and infers the
// log format based on the stdout device. If it looks like a terminal session, then it uses the console format.
func DefaultConfig() Config {
	format := "json"
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		format = "console" // looks like terminal session, use console logging
	}

	return Config{
		Level:  "info",
		Format: format,
	}
}

// Config encapsulates the configurable elements of the logger.
type Config struct {
	Level  string `json:"level"  usage:"adjust the verbosity of the logs" default:"info"`
	Format string `json:"format" usage:"configure the format of the logs" default:"json"`
}

// Setup creates a logger given the provided configuration.
func Setup(ctx context.Context, cfg Config) context.Context {
	level := zapcore.InfoLevel
	if cfg.Level != "" {
		err := (&level).Set(cfg.Level)
		if err != nil {
			panic(err)
		}
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.Level.SetLevel(level)
	zapConfig.Encoding = cfg.Format
	zapConfig.Sampling = nil // don't sample

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	return ToContext(ctx, logger)
}
