package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fabien-marty/slog-helpers/pkg/external"
)

// callback is called for each message logged
//
// attrs is a list of key/value pairs (as flattened strings)
func callback(t time.Time, level slog.Level, message string, attrs []external.StringifiedAttr) error {
	fmt.Println("time   :", t.Format(time.RFC3339))
	fmt.Println("level  :", level)
	fmt.Println("message:", message)
	for _, attr := range attrs {
		fmt.Printf("attr   : %s => %s\n", attr.Key, attr.Value)
	}
	return nil
}

func main() {
	// Create a new external handler
	handler := external.New(&external.Options{
		StringifiedCallback: callback, // use the simplified callback form (2 other forms are available)
	})

	// Create a logger with this handler
	logger := slog.New(handler)

	// Create a sub-logger with some default group/key
	logger = logger.With(slog.Group("common", slog.String("rootkey", "rootvalue")))

	// Log a message
	logger.Info("this is an info message", slog.String("key", "value"))
}
