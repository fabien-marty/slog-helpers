package slogh

import (
	"log/slog"
)

// Note: some MIT Code stolen from https://github.com/cappuccinotm/slogx

type payload struct {
	group  string
	attrs  []slog.Attr
	parent *payload
}

type accumulator struct {
	last *payload
}

// newAccumulator is a base struct for building slog.Handler that accumulates
// attributes and groups and returns them when calling assemble.
func newAccumulator() *accumulator {
	return &accumulator{
		last: &payload{},
	}
}

// WithAttrs returns a new accumulator with the given attributes.
func (a *accumulator) withAttrs(attrs []slog.Attr) *accumulator {
	acc := *a // shallow copy
	acc.last.attrs = append(acc.last.attrs, attrs...)
	return &acc
}

// WithGroup returns a new accumulator with the given group.
func (a *accumulator) withGroup(group string) *accumulator {
	acc := *a // shallow copy
	acc.last = &payload{group: group, parent: acc.last}
	return &acc
}

func (a *accumulator) assemble() (attrs []slog.Attr) {
	for p := a.last; p != nil; p = p.parent {
		attrs = append(p.attrs, attrs...)
		if p.group != "" {
			attrs = []slog.Attr{slog.Group(p.group, listAny(attrs)...)}
		}
	}
	return attrs
}

func (a *accumulator) withRecordAttrs(rec slog.Record) *accumulator {
	attrs := getAttrsFromRecord(rec)
	return a.withAttrs(attrs)
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
