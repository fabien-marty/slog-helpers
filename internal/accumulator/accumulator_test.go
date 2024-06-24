package accumulator

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAssemble(t *testing.T) {
	a := New()
	a = a.WithAttrs([]slog.Attr{slog.String("rootkey", "rootvalue")})
	a = a.WithGroup("group1")
	a = a.WithAttrs([]slog.Attr{slog.String("group1key", "group1value")})
	a = a.WithGroup("group2")
	a = a.WithAttrs([]slog.Attr{slog.String("group2key", "group2value")})
	attrs := a.Assemble()
	assert.Equal(t, 2, len(attrs))
	assert.Equal(t, "rootkey", attrs[0].Key)
	assert.Equal(t, "rootvalue", attrs[0].Value.String())
	group1 := attrs[1]
	assert.Equal(t, slog.KindGroup, group1.Value.Kind())
	assert.Equal(t, "group1", group1.Key)
	group1Attrs := group1.Value.Group()
	assert.Equal(t, 2, len(group1Attrs))
	assert.Equal(t, "group1key", group1Attrs[0].Key)
	assert.Equal(t, "group1value", group1Attrs[0].Value.String())
	group2 := group1Attrs[1]
	assert.Equal(t, slog.KindGroup, group2.Value.Kind())
	assert.Equal(t, "group2", group2.Key)
	group2Attrs := group2.Value.Group()
	assert.Equal(t, 1, len(group2Attrs))
	assert.Equal(t, "group2key", group2Attrs[0].Key)
	assert.Equal(t, "group2value", group2Attrs[0].Value.String())
}

func TestAssembleWithRecordAttrs(t *testing.T) {
	a := New()
	a = a.WithAttrs([]slog.Attr{slog.Int("foo", 123)}).WithGroup("group").WithAttrs([]slog.Attr{slog.String("foo2", "bar2")})
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "hello world", 0)
	record.AddAttrs(slog.String("foo3", "bar3"), slog.Group("zzz", slog.String("aaa", "bbb")))
	attrs := a.AssembleWithRecordAttrs(record)
	assert.Equal(t, 2, len(attrs))
	assert.Equal(t, "foo=123", attrs[0].String())
	assert.Equal(t, "group=[foo2=bar2 foo3=bar3 zzz=[aaa=bbb]]", attrs[1].String())
}

func TestRealCopy(t *testing.T) {
	a := New()
	b := a.WithAttrs([]slog.Attr{slog.String("foo", "bar")})
	attrs := b.Assemble()
	assert.Equal(t, 1, len(attrs))
	attrs = a.Assemble()
	assert.Equal(t, 0, len(attrs))
}
