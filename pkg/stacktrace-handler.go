package slogh

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/ztrue/tracerr"
)

// KeyNameForModeAddAttrDefault is the default key name for the StackTraceHandlerModeAddAttr mode.
const KeyNameForModeAddAttrDefault = "stacktrace"

// StackTraceHandlerMode is an enumeration type that defines the possible modes of the StackTraceHandler.
type StackTraceHandlerMode string

// StackTraceHandlerModeNothing is a mode that does nothing.
const StackTraceHandlerModeNothing StackTraceHandlerMode = "nothing"

// StackTraceHandlerModeAddAttr is a mode that adds an attribute to the record.
//
// The key name for the attribute can be overriden in KeyNameForModeAddAttr key in StackTraceHandlerOptions.
// The default value is defined in the KeyNameForModeAddAttrDefault constant.
const StackTraceHandlerModeAddAttr StackTraceHandlerMode = "add-attr"

// StackTraceHandlerModePrint is a mode that prints the stack trace to the output.
const StackTraceHandlerModePrint StackTraceHandlerMode = "print"

// StackTraceHandlerModePrintWithColors is a mode that prints the stack trace to the output with colors.
//
// Note that this option is only interesting when used with AddSource=true in the HandlerOptions.
const StackTraceHandlerModePrintWithColors StackTraceHandlerMode = "print-colors"

// The default mode of the StackTraceHandler.
const StackTraceHandlerModeDefault = StackTraceHandlerModePrint

var stackMutex sync.Mutex

// StackTraceHandlerOptions is a struct that contains the options for the StackTraceHandler.
type StackTraceHandlerOptions struct {
	slog.HandlerOptions
	Mode                  StackTraceHandlerMode // The mode of the StackTraceHandler.
	KeyNameForModeAddAttr string                // The key name for the attribute in StackTraceHandlerModeAddAttr mode.
	WriterForPrint        io.Writer             // The writer to use for ModePrint and ModePrintWithColors (default to stderr).
}

// StackTraceHandler is a slog handler that adds a stack trace to the record (add attribute or print/write).
//
// The stack trace is added/printed only if the StackTraceEnabled method returns true.
// The default behavior is to add the stack trace for records with a level greater or equal to slog.LevelError.
//
// The stack is added as an attribute if the Mode is StackTraceHandlerOptions is StackTraceHandlerModeAddAttr (great for JSON format for example).
// The stack can be dumped in a writer (default to stderr) if the Mode is StackTraceHandlerModePrint or StackTraceHandlerModePrintWithColors.
//
// Full example:
//
//	handler := slog.NewTextHandler(os.Stderr)
//	stackHandler := NewStackTraceHandler(handler, &StackTraceHandlerOptions{
//		Mode: StackTraceHandlerModePrintWithColors,
//	})
//	logger := slog.New(stackHandler)
//	logger.Info("no stack trace")
//	logger.Error("this is an error, let's print a stack trace")
type StackTraceHandler struct {
	slog.Handler
	opts *StackTraceHandlerOptions
}

// NewStackTraceHandler creates a new StackTraceHandler.
func NewStackTraceHandler(originalHandler slog.Handler, options *StackTraceHandlerOptions) slog.Handler {
	if options.WriterForPrint == nil {
		options.WriterForPrint = os.Stderr
	}
	if options.Mode == "" {
		options.Mode = StackTraceHandlerModeDefault
	}
	return &StackTraceHandler{
		Handler: originalHandler,
		opts:    options,
	}
}

// StackTraceEnabled returns true if the stack trace must be added/printed.
//
// The default behavior is to add the stack trace for records with a level greater or equal to slog.LevelError.
// You can override this method to customize the behavior.
//
// Important note: the behavior of this
func (sd *StackTraceHandler) StackTraceEnabled(context context.Context, record slog.Record) bool {
	return record.Level >= slog.LevelError
}

func (sd *StackTraceHandler) beforeHandle(record *slog.Record) error {
	switch sd.opts.Mode {
	case StackTraceHandlerModeAddAttr:
		fakeErr := tracerr.Wrap(errors.New("stack trace"))
		keyName := sd.opts.KeyNameForModeAddAttr
		if keyName == "" {
			keyName = KeyNameForModeAddAttrDefault
		}
		record.AddAttrs(slog.String(keyName, tracerr.Sprint(fakeErr)))
	}
	return nil
}

func (sd *StackTraceHandler) afterHandle(slog.Record) error {
	var str string
	var err error
	switch sd.opts.Mode {
	case StackTraceHandlerModePrint:
		fakeErr := tracerr.Wrap(errors.New(""))
		if sd.opts.AddSource {
			str = tracerr.SprintSource(fakeErr)
		} else {
			str = tracerr.Sprint(fakeErr)
		}
		stackMutex.Lock()
		defer stackMutex.Unlock()
		_, err = sd.opts.WriterForPrint.Write([]byte("error log level detected, let's print a stack trace" + "\n" + str + "\n"))
	case StackTraceHandlerModePrintWithColors:
		fakeErr := tracerr.Wrap(errors.New(""))
		if sd.opts.AddSource {
			str = tracerr.SprintSourceColor(fakeErr)
		} else {
			str = tracerr.Sprint(fakeErr)
		}
		stackMutex.Lock()
		defer stackMutex.Unlock()
		_, err = sd.opts.WriterForPrint.Write([]byte(ansiRedBackground + ansiWhite + "error log level detected, let's print a stack trace" + ansiReset + "\n" + str + "\n"))

	}
	return err
}

// Handle forwards the call to the original handler (see constructor) and adds/prints the stack trace if needed.
func (sd *StackTraceHandler) Handle(context context.Context, record slog.Record) error {
	var err error
	stackTraceEnabled := sd.StackTraceEnabled(context, record)
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
