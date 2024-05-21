package main

import (
	"log/slog"

	slogh "github.com/fabien-marty/slog-helpers/pkg"
)

func main() {
	logger := slogh.GetLogger(slogh.WithLevel(slog.LevelDebug), slogh.WithLogFormat(slogh.LogFormatTextHuman), slogh.WithStackTrace(true), slogh.WithColors(true))
	logger2 := logger.With("fab", "ien").WithGroup("xxxx").With("x", 4).WithGroup("yyyy").With("abc", "def")
	logger2.Warn("foo", slog.String("coucou", "foo"), slog.Group("zzz", slog.String("aaa", "bbb")))
	logger2.Debug("this is a debug message")
	logger.Info("this is an info message")
	newLogger := logger2.With(slog.Group("anothergroup", slog.String("key", "value"), slog.Group("anotherinsidegroup", slog.String("key", "value"))))
	newLogger.Info("coucou")
	foo(logger)
}

func foo(logger *slog.Logger) {
	bar(logger)
}

func bar(logger *slog.Logger) {
	logger.Error("this is an error")
}
