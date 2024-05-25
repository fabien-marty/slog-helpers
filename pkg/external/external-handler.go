package external

import (
	"context"
	"log/slog"
	"time"

	"github.com/fabien-marty/slog-helpers/internal/accumulator"
)

var _ slog.Handler = &Handler{}

// Callback is a function that handles nearly untouched slog log records.
type Callback func(time time.Time, level slog.Level, message string, attrs []slog.Attr) error

// FlattenedAttrsCallback is a function that handles slog log records with flattened attributes (no group, prefixed keys with group names).
type FlattenedAttrsCallback func(time time.Time, level slog.Level, message string, attrs []FlattenedAttr) error

// StringifiedAttrsCallback is a function that handles slog log records with stringified and flattened attributes (no group, prefixed keys with group names, values resolved as strings).
type StringifiedAttrsCallback func(time time.Time, level slog.Level, message string, attrs []StringifiedAttr) error

// Options is a struct that contains the options for the ExternalHandler.
type Options struct {
	slog.HandlerOptions
	Callback            Callback                 // If not nil, this callback will be used to handle the log records.
	FlattenedCallback   FlattenedAttrsCallback   // If not nil, this callback (with flattened attributes) will be used to handle the log records.
	StringifiedCallback StringifiedAttrsCallback // If not nil, this callback (with stringified and flattened attributes) will be used to handle the log records.
}

// Handler is an opaque type that implements the slog.Handler interface.
type Handler struct {
	*accumulator.Accumulator
	opts *Options
}

// New creates a new ExternalHandler.
func New(opts *Options) *Handler {
	return &Handler{
		Accumulator: accumulator.New(),
		opts:        opts,
	}
}

func (eh *Handler) Enabled(context context.Context, level slog.Level) bool {
	enabledLevel := slog.LevelInfo
	if eh.opts.HandlerOptions.Level != nil {
		enabledLevel = eh.opts.HandlerOptions.Level.Level()
	}
	return level >= enabledLevel
}

func (eh *Handler) WithGroup(group string) slog.Handler {
	eh.Accumulator = eh.Accumulator.WithGroup(group)
	return eh
}

func (eh *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	eh.Accumulator = eh.Accumulator.WithAttrs(attrs)
	return eh
}

func (eh *Handler) Handle(context context.Context, record slog.Record) error {
	var attrs []slog.Attr = eh.Accumulator.AssembleWithRecordAttrs(record)
	if eh.opts.Callback != nil {
		return eh.opts.Callback(record.Time, record.Level, record.Message, attrs)
	}
	if eh.opts.FlattenedCallback != nil {
		fattrs := newFlattenedAttrs(attrs, "")
		return eh.opts.FlattenedCallback(record.Time, record.Level, record.Message, fattrs)
	}
	if eh.opts.StringifiedCallback != nil {
		sattrs := newStringifiedAttrs(newFlattenedAttrs(attrs, ""))
		return eh.opts.StringifiedCallback(record.Time, record.Level, record.Message, sattrs)
	}
	return nil
}
