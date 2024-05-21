package slogh

import (
	"context"
	"log/slog"
	"time"
)

var _ slog.Handler = &ExternalHandler{}

// ExternalHandlerFunction is a function that handles nearly untouched slog log records.
type ExternalHandlerFunction func(time time.Time, level slog.Level, message string, attrs []slog.Attr) error

// ExternalHandlerFlattenedAttrsFunction is a function that handles slog log records with flattened attributes (no group, prefixed keys with group names).
type ExternalHandlerFlattenedAttrsFunction func(time time.Time, level slog.Level, message string, attrs []FlattenedAttr) error

// ExternalHandlerStringifiedAttrsFunction is a function that handles slog log records with stringified and flattened attributes (no group, prefixed keys with group names, values resolved as strings).
type ExternalHandlerStringifiedAttrsFunction func(time time.Time, level slog.Level, message string, attrs []StringifiedAttr) error

// ExternalHandlerOptions is a struct that contains the options for the ExternalHandler.
type ExternalHandlerOptions struct {
	slog.HandlerOptions
	Callback            ExternalHandlerFunction                 // If not nil, this callback will be used to handle the log records.
	FlattenedCallback   ExternalHandlerFlattenedAttrsFunction   // If not nil, this callback (with flattened attributes) will be used to handle the log records.
	StringifiedCallback ExternalHandlerStringifiedAttrsFunction // If not nil, this callback (with stringified and flattened attributes) will be used to handle the log records.
}

// ExternalHandler is an opaque type that implements the slog.Handler interface.
type ExternalHandler struct {
	*accumulator
	opts *ExternalHandlerOptions
}

// NewExternalHandler creates a new ExternalHandler.
func NewExternalHandler(opts *ExternalHandlerOptions) *ExternalHandler {
	return &ExternalHandler{
		accumulator: newAccumulator(),
		opts:        opts,
	}
}

func (eh *ExternalHandler) Enabled(context context.Context, level slog.Level) bool {
	enabledLevel := slog.LevelInfo
	if eh.opts.HandlerOptions.Level != nil {
		enabledLevel = eh.opts.HandlerOptions.Level.Level()
	}
	return level >= enabledLevel
}

func (eh *ExternalHandler) WithGroup(group string) slog.Handler {
	eh.accumulator = eh.withGroup(group)
	return eh
}

func (eh *ExternalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	eh.accumulator = eh.withAttrs(attrs)
	return eh
}

func (eh *ExternalHandler) Handle(context context.Context, record slog.Record) error {
	var attrs []slog.Attr = eh.accumulator.withRecordAttrs(record).assemble()
	if eh.opts.Callback != nil {
		return eh.opts.Callback(record.Time, record.Level, record.Message, attrs)
	}
	if eh.opts.FlattenedCallback != nil {
		fattrs := NewFlattenedAttrs(attrs, "")
		return eh.opts.FlattenedCallback(record.Time, record.Level, record.Message, fattrs)
	}
	if eh.opts.StringifiedCallback != nil {
		sattrs := NewStringifiedAttrs(NewFlattenedAttrs(attrs, ""))
		return eh.opts.StringifiedCallback(record.Time, record.Level, record.Message, sattrs)
	}
	return nil
}
