package accumulator

import (
	"log/slog"
	"slices"
)

// Note: some inspiration and some MIT Code stolen from https://github.com/cappuccinotm/slogx

type payload struct {
	group  string
	attrs  []slog.Attr
	parent *payload
}

type Accumulator struct {
	last *payload
}

func New() *Accumulator {
	return &Accumulator{
		last: &payload{},
	}
}

func (a *Accumulator) Clone() *Accumulator {
	return &Accumulator{last: &payload{
		group:  a.last.group,
		parent: a.last.parent,
		attrs:  slices.Clone(a.last.attrs),
	}}
}

// WithAttrs returns a new accumulator with the given attributes.
func (a *Accumulator) WithAttrs(attrs []slog.Attr) *Accumulator {
	acc := a.Clone()
	acc.last.attrs = append(acc.last.attrs, attrs...)
	return acc
}

// WithGroup returns a new accumulator with the given group.
func (a *Accumulator) WithGroup(group string) *Accumulator {
	acc := a.Clone()
	acc.last = &payload{group: group, parent: acc.last}
	return acc
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

func (a *Accumulator) AssembleWithRecordAttrs(rec slog.Record) (attrs []slog.Attr) {
	return a.Clone().WithAttrs(getRecordAttrs(rec)).Assemble()
}

func listAny(attrs []slog.Attr) []any {
	list := make([]any, len(attrs))
	for i, a := range attrs {
		list[i] = a
	}
	return list
}

func getRecordAttrs(rec slog.Record) []slog.Attr {
	attrs := []slog.Attr{}
	rec.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})
	return attrs
}
