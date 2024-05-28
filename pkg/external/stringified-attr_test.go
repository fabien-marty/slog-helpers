package external

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStringifiedAttr(t *testing.T) {
	a := slog.Int("foo", 1234)
	b := slog.String("bar", "baz")
	a1 := newFlattenedAttr(a, "group")
	b1 := newFlattenedAttr(b, "")
	a2 := newStringifiedAttr(a1)
	b2 := newStringifiedAttr(b1)
	assert.Equal(t, "group.foo", a2.Key)
	assert.Equal(t, "1234", a2.Value)
	assert.Equal(t, "bar", b2.Key)
	assert.Equal(t, "baz", b2.Value)
	assert.Equal(t, "group.foo=1234", a2.String())
}

func TestNewStringifiedAttrs(t *testing.T) {
	a := slog.Int("foo", 1234)
	b := slog.String("bar", "baz")
	a1 := newFlattenedAttr(a, "group")
	b1 := newFlattenedAttr(b, "")
	attrs := []FlattenedAttr{a1, b1}
	res := newStringifiedAttrs(attrs)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "group.foo=1234", res[0].String())
	assert.Equal(t, "bar=baz", res[1].String())
}
