package stacktrace

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/fabien-marty/slog-helpers/internal/ansi"

	"github.com/fabien-marty/tracerr"
)

// KeyForStackTraceEnabledDefault is the key of a boolean attribute to enable the stack trace for log level < MinimalLevelForStackTraceEnabledEnabled
//
// This behavior can be completely change by overriding the StackTraceEnabled method.
//
// Note: in the default behavior, this attribute will be automatically removed by this handler.
const KeyForStackTraceEnabledDefault = "add-stacktrace"

// MinimalLevelForStackTraceEnabledEnabled is the minimal level for which the stack trace is automatically enabled.
//
// This behavior can be completely change by overriding the StackTraceEnabled method.
const MinimalLevelForStackTraceEnabledEnabledDefault = slog.LevelError

// Mode is an enumeration type that defines the possible modes of the StackTraceHandler.
type Mode string

// ModeNothing is a mode that does nothing.
const ModeNothing Mode = "nothing"

// ModeAddAttr is a mode that adds an attribute to the record.
//
// The key name for the attribute can be overriden in KeyNameForModeAddAttr key in StackTraceHandlerOptions.
// The default value is defined in the KeyNameForModeAddAttrDefault constant.
const ModeAddAttr Mode = "add-attr"

// ModePrint is a mode that prints the stack trace to the output.
const ModePrint Mode = "print"

// ModePrintWithColors is a mode that prints the stack trace to the output with colors.
//
// Note that this option is only interesting when used with AddSource=true in the HandlerOptions.
const ModePrintWithColors Mode = "print-colors"

// The default mode of the StackTraceHandler.
const ModeDefault = ModePrint

// KeyNameForModeAddAttrDefault is the default key name for the ModeAddAttr mode.
const KeyNameForModeAddAttrDefault = "stacktrace"

var mutex sync.Mutex

func init() {
	// not great to do this tuning here but tracerr API could be better IMHO
	tracerr.DefaultIgnoreFirstFrames = 4
	tracerr.DefaultIgnoreLastFrames = 2
}

// Options is a struct that contains the options for the StackTraceHandler.
type Options struct {
	slog.HandlerOptions
	Mode                                    Mode        // The mode of the (stacktrace) Handler.
	KeyNameForModeAddAttr                   string      // The key name for the attribute in ModeAddAttr mode.
	KeyForStackTraceEnabled                 string      // The key of a boolean attribute to enable the stack trace
	MinimalLevelForStackTraceEnabledEnabled *slog.Level // The minimal level for which the stack trace is automatically enabled.
	WriterForPrint                          io.Writer   // The writer to use for ModePrint and ModePrintWithColors (default to stderr).
}

// Handler is a slog handler that adds a stack trace to the record (add attribute or print/write).
type Handler struct {
	slog.Handler
	opts *Options
}

// New creates a new Handler.
func New(originalHandler slog.Handler, options *Options) slog.Handler {
	if options.WriterForPrint == nil {
		options.WriterForPrint = os.Stderr
	}
	if options.Mode == "" {
		options.Mode = ModeDefault
	}
	if options.KeyForStackTraceEnabled == "" {
		options.KeyForStackTraceEnabled = KeyForStackTraceEnabledDefault
	}
	if options.MinimalLevelForStackTraceEnabledEnabled == nil {
		defaultLevel := MinimalLevelForStackTraceEnabledEnabledDefault
		options.MinimalLevelForStackTraceEnabledEnabled = &defaultLevel
	}
	return &Handler{
		Handler: originalHandler,
		opts:    options,
	}
}

func (sd *Handler) WithGroup(name string) slog.Handler {
	return New(sd.Handler.WithGroup(name), sd.opts)
}

func (sd *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return New(sd.Handler.WithAttrs(attrs), sd.opts)
}

// StackTraceEnabled returns true if the stack trace must be added/printed.
//
// The default behavior is to add the stack trace for records with a level greater or equal to slog.LevelError.
// You can override this method to customize the behavior.
//
// Important note: the behavior of this
func (sd *Handler) StackTraceEnabled(context context.Context, record *slog.Record) bool {
	var stackStraceEnabled bool
	newRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == sd.opts.KeyForStackTraceEnabled {
			if attr.Value.String() == "true" {
				stackStraceEnabled = true
			}
			return true // we don't add this attribute to the new record
		}
		newRecord.AddAttrs(attr)
		return true
	})
	*record = newRecord
	if sd.opts.MinimalLevelForStackTraceEnabledEnabled != nil && record.Level >= *sd.opts.MinimalLevelForStackTraceEnabledEnabled {
		stackStraceEnabled = true
	}
	return stackStraceEnabled
}

func (sd *Handler) beforeHandle(record *slog.Record) error {
	switch sd.opts.Mode {
	case ModeAddAttr:
		fakeErr := tracerr.Wrap(errors.New("stack trace"))
		keyName := sd.opts.KeyNameForModeAddAttr
		if keyName == "" {
			keyName = KeyNameForModeAddAttrDefault
		}
		record.AddAttrs(slog.String(keyName, tracerr.Sprint(fakeErr)))
	}
	return nil
}

func (sd *Handler) afterHandle(slog.Record) error {
	var str string
	var err error
	switch sd.opts.Mode {
	case ModePrint:
		fakeErr := tracerr.Wrap(errors.New(""))
		if sd.opts.AddSource {
			str = tracerr.SprintSource(fakeErr)
		} else {
			str = tracerr.Sprint(fakeErr)
		}
		mutex.Lock()
		defer mutex.Unlock()
		_, err = sd.opts.WriterForPrint.Write([]byte("stacktrace enabled, let's print a stack trace" + "\n" + str + "\n"))
	case ModePrintWithColors:
		fakeErr := tracerr.Wrap(errors.New(""))
		if sd.opts.AddSource {
			str = tracerr.SprintSourceColor(fakeErr)
		} else {
			str = tracerr.Sprint(fakeErr)
		}
		mutex.Lock()
		defer mutex.Unlock()
		_, err = sd.opts.WriterForPrint.Write([]byte(ansi.RedBackground + ansi.White + "stacktrace enabled, let's print a stack trace" + ansi.Reset + "\n" + str + "\n"))

	}
	return err
}

// Handle forwards the call to the original handler (see constructor) and adds/prints the stack trace if needed.
func (sd *Handler) Handle(context context.Context, record slog.Record) error {
	var err error
	stackTraceEnabled := sd.StackTraceEnabled(context, &record)
	if stackTraceEnabled {
		err = sd.beforeHandle(&record)
		if err != nil {
			return err
		}
	}
	err = sd.Handler.Handle(context, record)
	if err != nil {
		return err
	}
	if stackTraceEnabled {
		err = sd.afterHandle(record)
		if err != nil {
			return err
		}
	}
	return err
}
