package main

import (
	"log/slog"

	"github.com/fabien-marty/slog-helpers/pkg/slogc"
)

func main() {

	// Create a new slog.Logger automatically configured from environment variables and given options
	logger := slogc.GetLogger(
		slogc.WithLevel(slog.LevelDebug),
		slogc.WithStackTrace(true),
		slogc.WithColors(true),
	)

	// Create a sub-logger with some default group/key
	logger = logger.With(slog.Group("common", slog.String("rootkey", "rootvalue")))

	// Log some messages
	logger.Debug("this is a debug message", slog.String("key", "value"))
	logger.Info("this is an info message")
	logger.Warn("this is a warning message", slog.Int("intkey", 123))
	funcToShowcaseTheStackTrace(logger)

}

func funcToShowcaseTheStackTrace(lgr *slog.Logger) {
	lgr.Warn("this is a warning but with a stack trace", slog.Bool("add-stacktrace", true))
}
