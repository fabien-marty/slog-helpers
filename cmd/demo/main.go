package main

import (
	"log/slog"

	"github.com/fabien-marty/slog-helpers/pkg/configurator"
)

func main() {
	logger := configurator.GetLogger(
		configurator.WithLevel(slog.LevelDebug),
		configurator.WithLogFormat(configurator.LogFormatTextHuman),
		configurator.WithStackTrace(true),
		configurator.WithColors(true),
		configurator.WithStackTrace(true),
	)
	logger = logger.With("rootkey", "rootvalue")
	logger.Debug("this is a debug message", slog.String("key", "value"))
	logger.Info("this is an info message")
	anotherLogger := logger.WithGroup("group")
	anotherLogger.Warn("this is a warning message", slog.Int("intkey", 123))
	anotherFunction(anotherLogger) // log an error through another function to showcase the stacktrace
}

func anotherFunction(lgr *slog.Logger) {
	lgr.Error("this is an error with an automatic stackstrace")
}
