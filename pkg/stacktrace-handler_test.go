package slogh

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type record map[string]any

func TestStackTraceHandlerAddAttr(t *testing.T) {
	buffer := getBuffer()
	defer putBuffer(buffer)
	jsonHandler := slog.NewJSONHandler(buffer, &slog.HandlerOptions{})
	h := NewStackTraceHandler(jsonHandler, &StackTraceHandlerOptions{
		Mode: StackTraceHandlerModeAddAttr,
	})
	logger := slog.New(h)
	logger.Info("hello world")
	json1 := buffer.Bytes()
	r := record{}
	err := json.Unmarshal(json1, &r)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", r["msg"])
	assert.Equal(t, "INFO", r["level"])
	assert.Nil(t, r["stacktrace"])
	buffer.Reset()
	logger.Error("hello error")
	json2 := buffer.Bytes()
	r2 := record{}
	err = json.Unmarshal(json2, &r2)
	assert.NoError(t, err)
	assert.Equal(t, "hello error", r2["msg"])
	assert.Equal(t, "ERROR", r2["level"])
	stacktrace, ok := r2["stacktrace"]
	assert.True(t, ok)
	sstracktrace, ok := stacktrace.(string)
	assert.True(t, ok)
	assert.Greater(t, len(sstracktrace), 100)
}

func TestStackTraceHandlerPrint(t *testing.T) {
	buffer := getBuffer()
	defer putBuffer(buffer)
	eh := NewExternalHandler(&ExternalHandlerOptions{
		StringifiedCallback: func(time time.Time, level slog.Level, message string, attrs []StringifiedAttr) error {
			buffer.WriteString(level.String() + "/" + message)
			tmp := []string{}
			for _, attr := range attrs {
				tmp = append(tmp, attr.String())
			}
			buffer.WriteString(strings.Join(tmp, ", ") + "\n")
			return nil
		},
	})
	h := NewStackTraceHandler(eh, &StackTraceHandlerOptions{
		Mode:           StackTraceHandlerModePrint,
		WriterForPrint: buffer,
	})
	logger := slog.New(h)
	logger.Info("hello world")
	assert.Equal(t, "INFO/hello world\n", buffer.String())
	buffer.Reset()
	logger.Error("hello error")
	output := buffer.String()
	assert.True(t, strings.HasPrefix(output, "ERROR/hello error\nerror log level detected, let's print a stack trace\n"))
	assert.Greater(t, len(output), 500)
}

func TestStackTraceHandlerPrintWithColors(t *testing.T) {
	buffer := getBuffer()
	defer putBuffer(buffer)
	eh := NewExternalHandler(&ExternalHandlerOptions{
		StringifiedCallback: func(time time.Time, level slog.Level, message string, attrs []StringifiedAttr) error {
			buffer.WriteString(level.String() + "/" + message)
			tmp := []string{}
			for _, attr := range attrs {
				tmp = append(tmp, attr.String())
			}
			buffer.WriteString(strings.Join(tmp, ", ") + "\n")
			return nil
		},
	})
	h := NewStackTraceHandler(eh, &StackTraceHandlerOptions{
		Mode:           StackTraceHandlerModePrintWithColors,
		WriterForPrint: buffer,
	})
	logger := slog.New(h)
	logger.Info("hello world")
	assert.Equal(t, "INFO/hello world\n", buffer.String())
	buffer.Reset()
	logger.Error("hello error")
	output := buffer.String()
	fmt.Println(output)
	assert.True(t, strings.Contains(output, "error log level detected, let's print a stack trace"))
	assert.Greater(t, len(output), 500)
}
