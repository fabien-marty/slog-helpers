package slogc

import (
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/fabien-marty/slog-helpers/internal/bufferpool"
	"github.com/fabien-marty/slog-helpers/pkg/external"
	"github.com/fabien-marty/slog-helpers/pkg/stacktrace"
	"github.com/stretchr/testify/assert"
)

func replaceDigits(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return 'x'
		}
		return r
	}, s)
}

func TestGetLoggerDefault(t *testing.T) {
	buffer := bufferpool.Get()
	defer bufferpool.Put(buffer)
	l := GetLogger(WithDestinationWriter(buffer))
	l.Debug("foo")
	output := buffer.String()
	assert.Equal(t, 0, len(output)) // DEBUG is filtered by default
	l.Warn("foo", slog.String("bar", "baz"))
	output = replaceDigits(buffer.String())
	buffer.Reset()
	assert.Equal(t, "xxxx-xx-xxTxx:xx:xxZ [WARN ] foo {bar=baz}\n", output)
}

func TestGetLoggerMixedArgs(t *testing.T) {
	var decoded map[string]any = map[string]any{}
	buffer := bufferpool.Get()
	defer bufferpool.Put(buffer)
	l := GetLogger(WithLevel(slog.LevelDebug), WithDestinationWriter(buffer), WithLogFormat(LogFormatJson), WithStackTrace(true))
	l.Debug("foo")
	err := json.Unmarshal(buffer.Bytes(), &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "foo", decoded["msg"])
	assert.Equal(t, "DEBUG", decoded["level"])
	time := decoded["time"].(string)
	source := decoded["source"]
	sourceFile := source.(map[string]any)["file"].(string)
	assert.Greater(t, len(time), 10)
	assert.Greater(t, len(sourceFile), 10)
	assert.Nil(t, decoded[stacktrace.KeyNameForModeAddAttrDefault])
	buffer.Reset()
	l.Error("bar")
	err = json.Unmarshal(buffer.Bytes(), &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "bar", decoded["msg"])
	assert.Equal(t, "ERROR", decoded["level"])
	stack := decoded[stacktrace.KeyNameForModeAddAttrDefault]
	assert.Greater(t, len(stack.(string)), 10)
}

func TestGetLoggerExternal1(t *testing.T) {
	l := GetLogger(WithLogFormat(LogFormatExternal), WithExternalCallback(func(tim time.Time, l slog.Level, m string, attrs []slog.Attr) error {
		assert.Equal(t, slog.LevelWarn, l)
		assert.Equal(t, "foo", m)
		assert.Equal(t, 1, len(attrs))
		assert.Equal(t, "bar=baz", attrs[0].String())
		return nil
	}))
	l.Warn("foo", slog.String("bar", "baz"))
}

func TestGetLoggerExternal2(t *testing.T) {
	l := GetLogger(WithLogFormat(LogFormatExternal), WithExternalFlattenedAttrsCallback(func(tim time.Time, l slog.Level, m string, attrs []external.FlattenedAttr) error {
		assert.Equal(t, slog.LevelWarn, l)
		assert.Equal(t, "foo", m)
		assert.Equal(t, 1, len(attrs))
		assert.Equal(t, "bar=baz", attrs[0].String())
		return nil
	}))
	l.Warn("foo", slog.String("bar", "baz"))
}

func TestGetLoggerExternal3(t *testing.T) {
	l := GetLogger(WithLogFormat(LogFormatExternal), WithExternalStringifiedAttrsCallback(func(tim time.Time, l slog.Level, m string, attrs []external.StringifiedAttr) error {
		assert.Equal(t, slog.LevelWarn, l)
		assert.Equal(t, "foo", m)
		assert.Equal(t, 1, len(attrs))
		assert.Equal(t, "bar=baz", attrs[0].String())
		return nil
	}))
	l.Warn("foo", slog.String("bar", "baz"))
}
