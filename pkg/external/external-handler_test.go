package external

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewExternalHandlerStringified(t *testing.T) {
	logMessage := "hello world"
	callback := func(time time.Time, level slog.Level, message string, attrs []StringifiedAttr) error {
		assert.False(t, time.IsZero())
		assert.Equal(t, slog.LevelInfo, level)
		assert.Equal(t, logMessage, message)
		assert.Equal(t, 4, len(attrs))
		assert.Equal(t, "foo=123", attrs[0].String())
		assert.Equal(t, "group.foo2=bar2", attrs[1].String())
		assert.Equal(t, "group.foo3=bar3", attrs[2].String())
		assert.Equal(t, "group.zzz.aaa=bbb", attrs[3].String())
		return nil
	}
	h := New(&Options{
		StringifiedCallback: callback,
	})
	logger := slog.New(h).With(slog.Int("foo", 123)).WithGroup("group").With(slog.String("foo2", "bar2"))
	logger.Info(logMessage, slog.String("foo3", "bar3"), slog.Group("zzz", slog.String("aaa", "bbb")))
}

func TestNewExternalHandlerFlattened(t *testing.T) {
	logMessage := "hello world"
	callback := func(time time.Time, level slog.Level, message string, attrs []FlattenedAttr) error {
		assert.False(t, time.IsZero())
		assert.Equal(t, slog.LevelInfo, level)
		assert.Equal(t, logMessage, message)
		assert.Equal(t, 4, len(attrs))
		assert.Equal(t, "foo", attrs[0].Key)
		assert.Equal(t, int64(123), attrs[0].Value.Int64())
		assert.Equal(t, "group.foo2=bar2", attrs[1].String())
		assert.Equal(t, "group.foo3=bar3", attrs[2].String())
		assert.Equal(t, "group.zzz.aaa=bbb", attrs[3].String())
		return nil
	}
	h := New(&Options{
		FlattenedCallback: callback,
	})
	logger := slog.New(h).With(slog.Int("foo", 123)).WithGroup("group").With(slog.String("foo2", "bar2"))
	logger.Info(logMessage, slog.String("foo3", "bar3"), slog.Group("zzz", slog.String("aaa", "bbb")))
}

func TestNewExternalHandler(t *testing.T) {
	logMessage := "hello world"
	callback := func(time time.Time, level slog.Level, message string, attrs []slog.Attr) error {
		assert.False(t, time.IsZero())
		assert.Equal(t, slog.LevelInfo, level)
		assert.Equal(t, logMessage, message)
		assert.Equal(t, 2, len(attrs))
		assert.Equal(t, "foo=123", attrs[0].String())
		assert.Equal(t, "group=[foo2=bar2 foo3=bar3 zzz=[aaa=bbb]]", attrs[1].String())
		return nil
	}
	_ = callback
	h := New(&Options{
		Callback: callback,
	})
	logger := slog.New(h).With(slog.Int("foo", 123)).WithGroup("group").With(slog.String("foo2", "bar2"))
	logger.Info(logMessage, slog.String("foo3", "bar3"), slog.Group("zzz", slog.String("aaa", "bbb")))
}

func TestNewExternalHandlerWithoutCallback(t *testing.T) {
	h := New(&Options{
		HandlerOptions: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	})
	logger := slog.New(h)
	logger.Warn("this is a warning")
}

func TestNewExternalHandlerMultipleCalls(t *testing.T) {
	logMessage := "hello world"
	warnCalled := false
	callback := func(time time.Time, level slog.Level, message string, attrs []slog.Attr) error {
		switch level {
		case slog.LevelInfo:
			assert.Equal(t, 1, len(attrs))
		case slog.LevelWarn:
			assert.Equal(t, 0, len(attrs))
			warnCalled = true
		default:
			panic("unexpected level")
		}
		return nil
	}
	h := New(&Options{
		Callback: callback,
	})
	logger := slog.New(h)
	logger.Info(logMessage, slog.String("foo", "bar"))
	logger.Warn(logMessage)
	assert.True(t, warnCalled)
}
