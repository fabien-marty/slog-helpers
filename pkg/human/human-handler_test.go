package human

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/fabien-marty/slog-helpers/internal/ansi"
	"github.com/fabien-marty/slog-helpers/internal/bufferpool"
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

func TestNewHumanHandlerNoColor(t *testing.T) {
	buffer := bufferpool.GetBuffer()
	defer bufferpool.PutBuffer(buffer)
	h := NewHumanHandler(buffer, &HumandHandlerOptions{
		UseColors: false,
	})
	logger := slog.New(h)
	logger.Info("hello world")
	logger.Warn("hello warning", slog.String("foo", "bar"), slog.String("foofoo", "barbar"))
	res := buffer.String()
	lines := strings.Split(res, "\n")
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, 0, len(lines[2]))
	lines[0] = replaceDigits(lines[0])
	lines[1] = replaceDigits(lines[1])
	assert.Equal(t, "xxxx-xx-xxTxx:xx:xxZ [INFO ] hello world", lines[0])
	assert.Equal(t, "xxxx-xx-xxTxx:xx:xxZ [WARN ] hello warning {foo=bar foofoo=barbar}", lines[1])
}

func TestNewHumanHandlerColors(t *testing.T) {
	buffer := bufferpool.GetBuffer()
	h := NewHumanHandler(buffer, &HumandHandlerOptions{
		UseColors: true,
	})
	logger := slog.New(h)
	logger.Info("hello world")
	logger.Warn("hello warning", slog.String("foo", "bar"), slog.String("foo2", "bar2"))
	res := buffer.String()
	lines := strings.Split(res, "\n")
	assert.Equal(t, 4, len(lines))
	assert.Equal(t, 0, len(lines[3]))
	fmt.Println(lines[0][3:])
	lines[0] = replaceDigits(lines[0])
	lines[1] = replaceDigits(lines[1])
	assert.Equal(t, "▶ \x1b[xxmxxxx-xx-xxTxx:xx:xxZ\x1b[xm \x1b[xxm[INFO ]\x1b[xm \x1b[xmhello world\x1b[xm", lines[0])
	assert.Equal(t, "▶ \x1b[xxmxxxx-xx-xxTxx:xx:xxZ\x1b[xm \x1b[xxm[WARN ]\x1b[xm \x1b[xmhello warning\x1b[xm", lines[1])
	assert.Equal(t, "    ↳ \x1b[33mfoo\x1b[0m\x1b[1m=\x1b[0m\x1b[35mbar\x1b[0m \x1b[33mfoo2\x1b[0m\x1b[1m=\x1b[0m\x1b[35mbar2\x1b[0m", lines[2])
}

func TestLevelToStringNoColor(t *testing.T) {
	assert.Equal(t, "[DEBUG]", levelToStringNoColor(slog.LevelDebug))
	assert.Equal(t, "[INFO ]", levelToStringNoColor(slog.LevelInfo))
	assert.Equal(t, "[WARN ]", levelToStringNoColor(slog.LevelWarn))
	assert.Equal(t, "[ERROR]", levelToStringNoColor(slog.LevelError))
	assert.Equal(t, "[?????]", levelToStringNoColor(slog.Level(42)))
}

func TestLevelToString(t *testing.T) {
	assert.Equal(t, ansi.Green+"[DEBUG]"+ansi.Reset, levelToString(slog.LevelDebug))
	assert.Equal(t, ansi.Blue+"[INFO ]"+ansi.Reset, levelToString(slog.LevelInfo))
	assert.Equal(t, ansi.Red+"[WARN ]"+ansi.Reset, levelToString(slog.LevelWarn))
	assert.Equal(t, ansi.RedBackground+ansi.White+"[ERROR]"+ansi.Reset, levelToString(slog.LevelError))
	assert.Equal(t, ansi.Cyan+"[?????]"+ansi.Reset, levelToString(slog.Level(42)))
}
