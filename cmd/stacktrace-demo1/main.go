package main

import (
	"log/slog"
	"os"

	"github.com/fabien-marty/slog-helpers/pkg/stacktrace"
)

func main() {
	// Create a first handler
	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{})

	// Create a new StackTrace handler by decorating the first one
	stackTraceHandler := stacktrace.New(textHandler, &stacktrace.Options{
		HandlerOptions: slog.HandlerOptions{
			AddSource: true,
		},
		Mode: stacktrace.ModePrintWithColors,
	})

	// Create a logger with the StackTrace handler
	logger := slog.New(stackTraceHandler)

	logger.Warn("this is a standard warning")
	funcToShowcaseTheStackTrace(logger)
}

func funcToShowcaseTheStackTrace(logger *slog.Logger) {
	logger.Warn("this is a warning but with a stack trace",
		slog.Bool("add-stacktrace", true), // force stacktrace add/dump for this message
	)
}
