package configurator

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

// DefaultLogLevelEnvVar is the default environment variable used to define the default log level.
//
// The default value "LOG_LEVEL" can be overridden with SetLogLevelEnvVar.
const DefaultLogLevelEnvVar = "LOG_LEVEL"

// DefaultLogLevel is the default log level.
const DefaultLogLevel = slog.LevelInfo

var logLevelEnvVarMutex = sync.RWMutex{}
var logLevelEnvVar = DefaultLogLevelEnvVar

// SetLogLevelEnvVar sets the environment variable used to define the default log level.
func SetLogLevelEnvVar(envVar string) {
	logLevelEnvVarMutex.Lock()
	defer logLevelEnvVarMutex.Unlock()
	logLevelEnvVar = envVar
}

// GetLogLevelFromString returns the log level from a string.
//
// The log level is case insensitive. If the string is not recognized, the default log level is returned.
func GetLogLevelFromString(logLevel string) slog.Level {
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "CRIT", "CRITICAL", "ERR", "ERROR", "FATAL":
		return slog.LevelError
	}
	return DefaultLogLevel
}

// GetDefaultLogLevel returns the default log level.
//
// The default log level is defined by the environment variable LOG_LEVEL.
// If the environment variable is not set or empty, the default log level is INFO.
func GetDefaultLogLevel() slog.Level {
	logLevelEnvVarMutex.RLock()
	defer logLevelEnvVarMutex.RUnlock()
	logLevelAsString := strings.TrimSpace(os.Getenv(logLevelEnvVar))
	return GetLogLevelFromString(logLevelAsString)
}

// getLogLevel returns the log level from a pointer to a log level.
//
// If the pointer is nil, the default log level is returned.
func getLogLevel(level *slog.Level) slog.Level {
	if level == nil {
		return GetDefaultLogLevel()
	}
	return *level
}
