package slogh

import (
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLoggerDefault(t *testing.T) {
	buffer := getBuffer()
	defer putBuffer(buffer)
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
	buffer := getBuffer()
	defer putBuffer(buffer)
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
	assert.Nil(t, decoded[KeyNameForModeAddAttrDefault])
	buffer.Reset()
	l.Error("bar")
	err = json.Unmarshal(buffer.Bytes(), &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "bar", decoded["msg"])
	assert.Equal(t, "ERROR", decoded["level"])
	stack := decoded[KeyNameForModeAddAttrDefault]
	assert.Greater(t, len(stack.(string)), 10)
}
