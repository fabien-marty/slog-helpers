package slogh

import (
	"io"
	"log/slog"
	"sync"
	"time"
)

var humanMutex sync.Mutex

// HumanHandler is an opaque type that implements the slog.Handler interface.
type HumanHandler struct {
	ExternalHandler
}

// HumandHandlerOptions is a struct that contains the options for the HumanHandler.
type HumandHandlerOptions struct {
	slog.HandlerOptions
	UseColors bool // If true, use colors in the output.
}

// NewHumanHandler creates a new HumanHandler.
func NewHumanHandler(w io.Writer, opts *HumandHandlerOptions) *HumanHandler {
	var callback ExternalHandlerStringifiedAttrsFunction
	if opts.UseColors {
		callback = func(time time.Time, level slog.Level, message string, attrs []StringifiedAttr) error {
			return handleColor(w, time, level, message, attrs)
		}
	} else {
		callback = func(time time.Time, level slog.Level, message string, attrs []StringifiedAttr) error {
			return handleNoColor(w, time, level, message, attrs)
		}
	}
	return &HumanHandler{
		ExternalHandler: *NewExternalHandler(&ExternalHandlerOptions{
			HandlerOptions:      opts.HandlerOptions,
			StringifiedCallback: callback,
		}),
	}
}

func levelToStringNoColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "[DEBUG]"
	case slog.LevelInfo:
		return "[INFO ]"
	case slog.LevelWarn:
		return "[WARN ]"
	case slog.LevelError:
		return "[ERROR]"
	}
	return "[?????]"
}

func levelToString(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return ansiGreen + "[DEBUG]" + ansiReset
	case slog.LevelInfo:
		return ansiBlue + "[INFO ]" + ansiReset
	case slog.LevelWarn:
		return ansiRed + "[WARN ]" + ansiReset
	case slog.LevelError:
		return ansiRedBackground + ansiWhite + "[ERROR]" + ansiReset
	}
	return ansiCyan + "[?????]" + ansiReset
}

func handleColor(w io.Writer, time time.Time, level slog.Level, message string, attrs []StringifiedAttr) error {
	buffer := getBuffer()
	defer putBuffer(buffer)
	ascTime := "              "
	if !time.IsZero() {
		ascTime = time.UTC().Format("2006-01-02T15:04:05Z")
	}
	buffer.WriteString("▶ ")
	buffer.WriteString(ansiCyan)
	buffer.WriteString(ascTime)
	buffer.WriteString(ansiReset)
	buffer.WriteString(" ")
	buffer.WriteString(levelToString(level))
	buffer.WriteString(" ")
	buffer.WriteString(ansiBold)
	buffer.WriteString(message)
	buffer.WriteString(ansiReset)
	nAttr := 0
	if len(attrs) > 0 {
		buffer.WriteString("\n    ↳ ")
	}
	for _, attr := range attrs {
		buffer.WriteString(ansiYellow)
		buffer.WriteString(attr.Key)
		buffer.WriteString(ansiReset)
		buffer.WriteString(ansiBold)
		buffer.WriteString("=")
		buffer.WriteString(ansiReset)
		buffer.WriteString(ansiMagenta)
		buffer.WriteString(attr.Value)
		buffer.WriteString(ansiReset)
		nAttr++
		if nAttr < len(attrs) {
			buffer.WriteString(" ")
		}
	}
	buffer.WriteString("\n")
	humanMutex.Lock()
	defer humanMutex.Unlock()
	_, err := w.Write(buffer.Bytes())
	return err
}

func handleNoColor(w io.Writer, time time.Time, level slog.Level, message string, attrs []StringifiedAttr) error {
	buffer := getBuffer()
	defer putBuffer(buffer)
	ascTime := "              "
	if !time.IsZero() {
		ascTime = time.UTC().Format("2006-01-02T15:04:05Z")
	}
	buffer.WriteString(ascTime)
	buffer.WriteString(" ")
	buffer.WriteString(levelToStringNoColor(level))
	buffer.WriteString(" ")
	buffer.WriteString(message)
	nAttr := 0
	if len(attrs) > 0 {
		buffer.WriteString(" {")
	}
	for _, attr := range attrs {
		buffer.WriteString(attr.Key)
		buffer.WriteString("=")
		buffer.WriteString(attr.Value)
		nAttr++
		if nAttr < len(attrs) {
			buffer.WriteString(" ")
		} else {
			buffer.WriteString("}")
		}
	}
	buffer.WriteString("\n")
	humanMutex.Lock()
	defer humanMutex.Unlock()
	_, err := w.Write(buffer.Bytes())
	return err
}
