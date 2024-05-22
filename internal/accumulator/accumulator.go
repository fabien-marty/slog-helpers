package accumulator

import (
	"log/slog"
)

// Note: some MIT Code stolen from https://github.com/cappuccinotm/slogx

type payload struct {
	group  string
	attrs  []slog.Attr
	parent *payload
}

type Accumulator struct {
	last *payload
}

// New is a base struct for building slog.Handler that accumulates
// attributes and groups and returns them when calling assemble.
func New() *Accumulator {
	return &Accumulator{
		last: &payload{},
	}
}

// WithAttrs returns a new accumulator with the given attributes.
func (a *Accumulator) WithAttrs(attrs []slog.Attr) *Accumulator {
	acc := *a // shallow copy
	acc.last.attrs = append(acc.last.attrs, attrs...)
	return &acc
}

// WithGroup returns a new accumulator with the given group.
func (a *Accumulator) WithGroup(group string) *Accumulator {
	acc := *a // shallow copy
	acc.last = &payload{group: group, parent: acc.last}
	return &acc
}

func (a *Accumulator) Assemble() (attrs []slog.Attr) {
	for p := a.last; p != nil; p = p.parent {
		attrs = append(p.attrs, attrs...)
		if p.group != "" {
			attrs = []slog.Attr{slog.Group(p.group, listAny(attrs)...)}
		}
	}
	return attrs
}

func (a *Accumulator) WithRecordAttrs(rec slog.Record) *Accumulator {
	attrs := getAttrsFromRecord(rec)
	return a.WithAttrs(attrs)
}

func listAny(attrs []slog.Attr) []any {
	list := make([]any, len(attrs))
	for i, a := range attrs {
		list[i] = a
	}
	return list
}

// getAttrsFromRecord is an utility function to return all attributes from the given record.
func getAttrsFromRecord(rec slog.Record) []slog.Attr {
	var attrs []slog.Attr
	rec.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})
	return attrs
}
