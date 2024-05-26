// stacktrace.Handler is a slog handler that adds a stack trace to the record (add attribute or print/write).
//
// It does not output the log record itself, it decorates another slog.Handler given in the New() method.
//
// The stack trace is added/printed only if the StackTraceEnabled method returns true.
// The default behavior is to add the stack trace for records with:
//   - a level greater or equal to slog.LevelError
//   - (or) a boolean attribute add-stacktrace=true
//
// Note: in the default behavior, the attribute "add-stacktrace" will be automatically removed by this handler.
//
// The stacktrace can be added as an attribute to the record if the Mode is StackTraceHandlerOptions is ModeAddAttr (great for JSON format for example).
// The stacktrace can be dumped in a writer (default to stderr) if the Mode is ModePrint or ModePrintWithColors.
//
// Full example:
//
//	handler := slog.NewTextHandler(os.Stderr)
//	stackHandler := New(handler, &Options{
//		Mode: ModePrintWithColors,
//	})
//	logger := slog.New(stackHandler)
//	logger.Warn("warning with no stack trace")
//	logger.Error("this is an error, let's print a stack trace")
//	logger.Warn("this a warning but with a stack trace", slog.Bool("add-stacktrace", true))
package stacktrace
