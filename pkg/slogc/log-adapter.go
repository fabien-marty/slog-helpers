package slogc

import (
	"log"
	"log/slog"
	"time"

	"github.com/fabien-marty/slog-helpers/pkg/external"
)

func NewLogSlogAdapter(originalLogger *log.Logger) *slog.Logger {
	var callback external.StringifiedAttrsCallback = func(time time.Time, level slog.Level, message string, attrs []external.StringifiedAttr) error {
		return callback(originalLogger, time, level, message, attrs)
	}
	return GetLogger(WithLogFormat(LogFormatExternal), WithExternalStringifiedAttrsCallback(callback))
}

func callback(originalLogger *log.Logger, time time.Time, level slog.Level, message string, attrs []external.StringifiedAttr) error {
	// Do something with the log message
	ascTime := "              "
	if !time.IsZero() {
		ascTime = time.UTC().Format("2006-01-02T15:04:05Z")
	}
	extra := ""
	if len(attrs) > 0 {
		extra = " {"
		for i, attr := range attrs {
			if i == 0 {
				extra += attr.String()
			} else {
				extra += " " + attr.String()
			}
		}
		extra = "}"
	}
	originalLogger.Printf("%s %s: %s%s\n", ascTime, level.String(), message, extra)
	return nil
}
