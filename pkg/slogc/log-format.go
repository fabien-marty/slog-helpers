package slogc

import (
	"os"
	"strings"
	"sync"
)

type LogFormat string

// DefaultLogFormatEnvVar is the default environment variable used to define the default log format.
//
// The default value "LOG_FORMAT" can be overridden with SetLogFormatEnvVar.
const DefaultLogFormatEnvVar = "LOG_FORMAT"

// LogFormatTextHuman is the human readable console format.
const LogFormatTextHuman LogFormat = "text-human"

// LogFormatText is the basic/standard text format.
const LogFormatText LogFormat = "text"

// LogFormatJson is the basic/standard JSON format.
const LogFormatJson LogFormat = "json"

// LogFormatJsonGcp is the JSON format for Google Cloud Platform (GCP).
const LogFormatJsonGcp LogFormat = "json-gcp"

// LogFormatExternal is the external format (log records are not rendered by the logger but sent to an external handler)
const LogFormatExternal LogFormat = "external"

// DefaultLogFormat is the default log format.
const DefaultLogFormat = LogFormatTextHuman

var logFormatEnvVarMutex = sync.RWMutex{}
var logFormatEnvVar = DefaultLogFormatEnvVar

// SetLogFormatEnvVar sets the environment variable used to define the default log format.
func SetLogFormatEnvVar(envVar string) {
	logFormatEnvVarMutex.Lock()
	defer logFormatEnvVarMutex.Unlock()
	logFormatEnvVar = envVar
}

// GetLogFormatFromString returns the log format from a string.
//
// The log format is case insensitive. If the string is not recognized, the default log format is returned.
func GetLogFormatFromString(logLevel string) LogFormat {
	switch strings.ToLower(logLevel) {
	case "text-human":
		return LogFormatTextHuman
	case "text":
		return LogFormatText
	case "json":
		return LogFormatJson
	case "json-gcp", "gcp":
		return LogFormatJsonGcp
	case "external":
		return LogFormatExternal
	}
	return DefaultLogFormat
}

// GetDefaultLogFormat returns the default log format.
//
// The default log format is defined by the environment variable LOG_FORMAT.
// If the environment variable is not set or empty, the default log format is human-text.
func GetDefaultLogFormat() LogFormat {
	logFormatEnvVarMutex.RLock()
	defer logFormatEnvVarMutex.RUnlock()
	logFormatAsString := strings.TrimSpace(os.Getenv(logFormatEnvVar))
	return GetLogFormatFromString(logFormatAsString)
}

func getLogFormat(format *LogFormat) LogFormat {
	if format == nil {
		return GetDefaultLogFormat()
	}
	return *format
}
