package stacktrace

import (
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/fabien-marty/slog-helpers/internal/bufferpool"
	"github.com/fabien-marty/slog-helpers/pkg/external"
	"github.com/stretchr/testify/assert"
)

type record map[string]any

func TestStackTraceHandlerAddAttr(t *testing.T) {
	buffer := bufferpool.Get()
	defer bufferpool.Put(buffer)
	jsonHandler := slog.NewJSONHandler(buffer, &slog.HandlerOptions{})
	h := New(jsonHandler, &Options{
		Mode: ModeAddAttr,
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
	buffer := bufferpool.Get()
	defer bufferpool.Put(buffer)
	eh := external.New(&external.Options{
		StringifiedCallback: func(time time.Time, level slog.Level, message string, attrs []external.StringifiedAttr) error {
			buffer.WriteString(level.String() + "/" + message)
			tmp := []string{}
			for _, attr := range attrs {
				tmp = append(tmp, attr.String())
			}
			buffer.WriteString(strings.Join(tmp, ", ") + "\n")
			return nil
		},
	})
	h := New(eh, &Options{
		Mode:           ModePrint,
		WriterForPrint: buffer,
	})
	logger := slog.New(h)
	logger.Info("hello world")
	assert.Equal(t, "INFO/hello world\n", buffer.String())
	buffer.Reset()
	logger.Error("hello error")
	output := buffer.String()
	assert.True(t, strings.HasPrefix(output, "ERROR/hello error\nstacktrace enabled, let's print a stack trace\n"))
	assert.Greater(t, len(output), 100)
}

func TestStackTraceHandlerPrintWithColors(t *testing.T) {
	buffer := bufferpool.Get()
	defer bufferpool.Put(buffer)
	eh := external.New(&external.Options{
		StringifiedCallback: func(time time.Time, level slog.Level, message string, attrs []external.StringifiedAttr) error {
			buffer.WriteString(level.String() + "/" + message)
			tmp := []string{}
			for _, attr := range attrs {
				tmp = append(tmp, attr.String())
			}
			buffer.WriteString(strings.Join(tmp, ", ") + "\n")
			return nil
		},
	})
	h := New(eh, &Options{
		Mode:           ModePrintWithColors,
		WriterForPrint: buffer,
	})
	logger := slog.New(h)
	logger.Info("hello world")
	assert.Equal(t, "INFO/hello world\n", buffer.String())
	buffer.Reset()
	logger.Error("hello error")
	output := buffer.String()
	assert.True(t, strings.Contains(output, "stacktrace enabled, let's print a stack trace"))
	assert.Greater(t, len(output), 100)
}

func TestStackTraceHandlerWithAttr(t *testing.T) {
	buffer := bufferpool.Get()
	defer bufferpool.Put(buffer)
	jsonHandler := slog.NewJSONHandler(buffer, &slog.HandlerOptions{})
	h := New(jsonHandler, &Options{
		Mode: ModeAddAttr,
	})
	logger := slog.New(h)
	logger.Info("hello world", slog.Bool("add-stacktrace", true), slog.String("key", "value"), slog.String("key2", "value2"))
	json1 := buffer.Bytes()
	r := record{}
	err := json.Unmarshal(json1, &r)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", r["msg"])
	assert.Equal(t, "INFO", r["level"])
	assert.Equal(t, "value", r["key"])
	assert.Equal(t, "value2", r["key2"])
	_, ok := r["add-stacktrace"]
	assert.False(t, ok)
	stacktrace, ok := r["stacktrace"]
	assert.True(t, ok)
	sstracktrace, ok := stacktrace.(string)
	assert.True(t, ok)
	assert.Greater(t, len(sstracktrace), 100)
}
