package slogc

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/fabien-marty/slog-helpers/pkg/external"
	"github.com/fabien-marty/slog-helpers/pkg/human"
	"github.com/fabien-marty/slog-helpers/pkg/stacktrace"
	"github.com/mattn/go-isatty"
	"github.com/vlad-tokarev/sloggcp"
)

type loggerOptions struct {
	_level                           *slog.Level
	_destination                     *LogDestination
	_format                          *LogFormat
	_stackTrace                      *bool
	_colors                          *bool
	destinationWriter                io.Writer
	externalCallback                 external.Callback
	externalFlattenedAttrsCallback   external.FlattenedAttrsCallback
	externalStringifiedAttrsCallback external.StringifiedAttrsCallback
	level                            slog.Level
	destination                      LogDestination
	format                           LogFormat
	stackTrace                       bool
	addSource                        bool
	colors                           bool
}

// LoggerOption is a type that defines the options for the logger.
//
// You don't need to use it.
type LoggerOption func(options *loggerOptions) error

// WithLevel is an option that sets the level of the logger.
func WithLevel(level slog.Level) LoggerOption {
	return func(options *loggerOptions) error {
		options._level = &level
		return nil
	}
}

// WithDestination is an option that sets the destination of the logger.
//
// Note: you can also use WithDestinationWriter to set a custom writer.
func WithDestination(destination LogDestination) LoggerOption {
	return func(options *loggerOptions) error {
		options._destination = &destination
		return nil
	}
}

// WithDestinationWriter is an option that sets the writer of the logger.
//
// Note: it overrides the destination set by WithDestination.
func WithDestinationWriter(destinationWriter io.Writer) LoggerOption {
	return func(options *loggerOptions) error {
		options.destinationWriter = destinationWriter
		return nil
	}
}

// WithLogFormat is an option that sets the format of the logger.
func WithLogFormat(format LogFormat) LoggerOption {
	return func(options *loggerOptions) error {
		options._format = &format
		return nil
	}
}

// WithStackTrace is an option that sets if the logger should print or add stack traces.
func WithStackTrace(flag bool) LoggerOption {
	return func(options *loggerOptions) error {
		options._stackTrace = &flag
		return nil
	}
}

// WithColors is an option that sets if the logger should use colors.
//
// If not used, the use of colors is automatic (depending on the terminal connected to the logger destination).
func WithColors(flag bool) LoggerOption {
	return func(options *loggerOptions) error {
		options._colors = &flag
		return nil
	}
}

func WithExternalCallback(callback external.Callback) LoggerOption {
	return func(options *loggerOptions) error {
		options.externalCallback = callback
		return nil
	}
}

func WithExternalFlattenedAttrsCallback(callback external.FlattenedAttrsCallback) LoggerOption {
	return func(options *loggerOptions) error {
		options.externalFlattenedAttrsCallback = callback
		return nil
	}
}

func WithExternalStringifiedAttrsCallback(callback external.StringifiedAttrsCallback) LoggerOption {
	return func(options *loggerOptions) error {
		options.externalStringifiedAttrsCallback = callback
		return nil
	}
}

func completeOptions(options *loggerOptions) {
	options.level = getLogLevel(options._level)
	options.destination = getDestination(options._destination)
	if options.destinationWriter == nil {
		options.destinationWriter = options.destination.getFile()
	}
	options.format = getLogFormat(options._format)
	if options._stackTrace != nil {
		options.stackTrace = *options._stackTrace
	} else {
		options.stackTrace = (options.level == slog.LevelDebug) && (options.format == LogFormatTextHuman)
	}
	if options._colors != nil {
		options.colors = *options._colors
	} else {
		file, ok := options.destinationWriter.(*os.File)
		if !ok || file == nil {
			options.colors = false
		} else {
			options.colors = isatty.IsTerminal(file.Fd())
		}
	}
	options.addSource = (options.level == slog.LevelDebug)
	if options.externalCallback != nil || options.externalFlattenedAttrsCallback != nil || options.externalStringifiedAttrsCallback != nil {
		options.format = LogFormatExternal // if an external callback is set, the format is forced to external
	}
}

// GetLogger creates a new configured logger with the given options.
//
// Hint for your IDE: all LoggerOption functions starts with "With".
func GetLogger(opts ...LoggerOption) *slog.Logger {
	var options loggerOptions
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			panic(err)
		}
	}
	completeOptions(&options)
	standardHandlerOpts := slog.HandlerOptions{
		Level:     options.level,
		AddSource: options.addSource,
	}
	var handler slog.Handler
	switch options.format {
	case LogFormatTextHuman:
		handler = human.New(options.destinationWriter, &human.Options{
			HandlerOptions: standardHandlerOpts,
			UseColors:      options.colors,
		},
		)
	case LogFormatText:
		handler = slog.NewTextHandler(options.destinationWriter, &standardHandlerOpts)
	case LogFormatJson:
		handler = slog.NewJSONHandler(options.destinationWriter, &standardHandlerOpts)
	case LogFormatJsonGcp:
		standardHandlerOpts.ReplaceAttr = sloggcp.ReplaceAttr
		handler = slog.NewJSONHandler(options.destinationWriter, &standardHandlerOpts)
	case LogFormatExternal:
		if options.externalCallback != nil {
			handler = external.New(&external.Options{
				HandlerOptions: standardHandlerOpts,
				Callback:       options.externalCallback,
			})
		} else if options.externalFlattenedAttrsCallback != nil {
			handler = external.New(&external.Options{
				HandlerOptions:    standardHandlerOpts,
				FlattenedCallback: options.externalFlattenedAttrsCallback,
			})
		} else if options.externalStringifiedAttrsCallback != nil {
			handler = external.New(&external.Options{
				HandlerOptions:      standardHandlerOpts,
				StringifiedCallback: options.externalStringifiedAttrsCallback,
			})
		} else {
			panic("log format = external but no callback provided")
		}
	default:
		panic(fmt.Sprintf("unsupported log format: %s", options.format))
	}
	if options.stackTrace {
		var mode stacktrace.Mode
		switch options.format {
		case LogFormatJsonGcp, LogFormatJson:
			mode = stacktrace.ModeAddAttr
		case LogFormatTextHuman, LogFormatText:
			if options.colors {
				mode = stacktrace.ModePrintWithColors
			} else {
				mode = stacktrace.ModePrint
			}
		}
		handler = stacktrace.New(handler, &stacktrace.Options{
			Mode:           mode,
			HandlerOptions: standardHandlerOpts,
			WriterForPrint: options.destinationWriter,
		},
		)
	}
	logger := slog.New(handler)
	return logger
}

// SetDefaultLogger configures a new logger and sets it as the default logger to be returned by slog.Default() calls or used by slog.Info/Debug/Warning/Error calls.
//
// This is the same than a GetLogger call followed by a slog.SetDefault call.
// See GetLogger
func SetDefaultLogger(opts ...LoggerOption) {
	logger := GetLogger(opts...)
	slog.SetDefault(logger)
}
