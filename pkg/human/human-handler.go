package human

import (
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/fabien-marty/slog-helpers/internal/ansi"
	"github.com/fabien-marty/slog-helpers/internal/bufferpool"
	"github.com/fabien-marty/slog-helpers/pkg/external"
)

var _ slog.Handler = &Handler{}

var mutex sync.Mutex

// Handler is an opaque type that implements the slog.Handler interface.
type Handler struct {
	external.Handler
}

// Options is a struct that contains the options for the HumanHandler.
type Options struct {
	slog.HandlerOptions
	UseColors bool // If true, use colors in the output.
}

// New creates a new HumanHandler.
func New(w io.Writer, opts *Options) *Handler {
	var callback external.StringifiedAttrsCallback
	if opts.UseColors {
		callback = func(time time.Time, level slog.Level, message string, attrs []external.StringifiedAttr) error {
			return handleColor(w, time, level, message, attrs)
		}
	} else {
		callback = func(time time.Time, level slog.Level, message string, attrs []external.StringifiedAttr) error {
			return handleNoColor(w, time, level, message, attrs)
		}
	}
	return &Handler{
		Handler: *external.New(&external.Options{
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
		return ansi.Green + "[DEBUG]" + ansi.Reset
	case slog.LevelInfo:
		return ansi.Blue + "[INFO ]" + ansi.Reset
	case slog.LevelWarn:
		return ansi.Red + "[WARN ]" + ansi.Reset
	case slog.LevelError:
		return ansi.RedBackground + ansi.White + "[ERROR]" + ansi.Reset
	}
	return ansi.Cyan + "[?????]" + ansi.Reset
}

func handleColor(w io.Writer, time time.Time, level slog.Level, message string, attrs []external.StringifiedAttr) error {
	buffer := bufferpool.Get()
	defer bufferpool.Put(buffer)
	ascTime := "              "
	if !time.IsZero() {
		ascTime = time.UTC().Format("2006-01-02T15:04:05Z")
	}
	buffer.WriteString("▶ ")
	buffer.WriteString(ansi.Cyan)
	buffer.WriteString(ascTime)
	buffer.WriteString(ansi.Reset)
	buffer.WriteString(" ")
	buffer.WriteString(levelToString(level))
	buffer.WriteString(" ")
	buffer.WriteString(ansi.Bold)
	buffer.WriteString(message)
	buffer.WriteString(ansi.Reset)
	nAttr := 0
	if len(attrs) > 0 {
		buffer.WriteString("\n    ↳ ")
	}
	for _, attr := range attrs {
		buffer.WriteString(ansi.Yellow)
		buffer.WriteString(attr.Key)
		buffer.WriteString(ansi.Reset)
		buffer.WriteString(ansi.Bold)
		buffer.WriteString("=")
		buffer.WriteString(ansi.Reset)
		buffer.WriteString(ansi.Magenta)
		buffer.WriteString(attr.Value)
		buffer.WriteString(ansi.Reset)
		nAttr++
		if nAttr < len(attrs) {
			buffer.WriteString(" ")
		}
	}
	buffer.WriteString("\n")
	mutex.Lock()
	defer mutex.Unlock()
	_, err := w.Write(buffer.Bytes())
	return err
}

func handleNoColor(w io.Writer, time time.Time, level slog.Level, message string, attrs []external.StringifiedAttr) error {
	buffer := bufferpool.Get()
	defer bufferpool.Put(buffer)
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
	mutex.Lock()
	defer mutex.Unlock()
	_, err := w.Write(buffer.Bytes())
	return err
}
