package main

import (
	"log/slog"
	"os"

	"github.com/fabien-marty/slog-helpers/pkg/human"
)

func main() {
	// Create a new human handler
	humanHandler := human.New(os.Stderr, &human.Options{
		HandlerOptions: slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
		UseColors: true, // force the usage of colors (the default behavior is to detect if the output is a terminal)
	})

	// Create a logger with this handler
	logger := slog.New(humanHandler)

	// Create a sub-logger with some default group/key
	logger = logger.With(slog.Group("common", slog.String("rootkey", "rootvalue")))

	// Log some messages
	logger.Debug("this is a debug message", slog.String("key", "value"))
	logger.Info("this is an info message")
	logger.Warn("this is a warning message", slog.Int("intkey", 123))
	logger.Error("this is an error message")
}
