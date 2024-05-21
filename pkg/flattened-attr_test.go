package slogh

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFlattenedAttr(t *testing.T) {
	attr := slog.String("key", "value")
	prefixed := NewFlattenedAttr(attr, "")
	assert.Equal(t, "key", prefixed.Key)
	assert.Equal(t, "value", prefixed.Value.String())
	prefixed = NewFlattenedAttr(attr, "foo")
	assert.Equal(t, "foo.key", prefixed.Key)
	assert.Equal(t, "value", prefixed.Value.String())
}

func TestNewFlattenedAttrs(t *testing.T) {
	a1 := slog.String("rootkey", "rootvalue")
	a2 := slog.Group(
		"group1",
		slog.String("group1key", "group1value"),
		slog.Group("group2", slog.String("group2key", "group2value")),
	)
	attrs := []slog.Attr{a1, a2}
	fattrs := NewFlattenedAttrs(attrs, "")
	assert.Equal(t, 3, len(fattrs))
	assert.Equal(t, "rootkey", fattrs[0].Key)
	assert.Equal(t, "rootvalue", fattrs[0].Value.String())
	assert.Equal(t, "group1.group1key", fattrs[1].Key)
	assert.Equal(t, "group1value", fattrs[1].Value.String())
	assert.Equal(t, "group1.group2.group2key", fattrs[2].Key)
	assert.Equal(t, "group2value", fattrs[2].Value.String())
}
