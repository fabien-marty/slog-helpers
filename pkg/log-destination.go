package slogh

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// LogDestination represents the destination of the logs.
type LogDestination string

// DefaultLogDestinationEnvVar is the default environment variable used to define the default log destination.
//
// The default value "LOG_DESTINATION" can be overridden with SetLogDestinationEnvVar.
const DefaultLogDestinationEnvVar = "LOG_DESTINATION"

// LogDestinationStdout is the standard output.
var LogDestinationStdout LogDestination = "stdout"

// LogDestinationStderr is the standard error.
var LogDestinationStderr LogDestination = "stderr"

// DefaultLogDestination is the default log destination.
var DefaultLogDestination = LogDestinationStderr

var logDestinationEnvVarMutex = sync.RWMutex{}
var logDestinationEnvVar = DefaultLogDestinationEnvVar

// SetLogDestinationEnvVar sets the environment variable used to define the default log destination.
func SetLogDestinationEnvVar(envVar string) {
	logDestinationEnvVarMutex.Lock()
	defer logDestinationEnvVarMutex.Unlock()
	logDestinationEnvVar = envVar
}

// GetLogDestinationFromString returns the log destination from a string.
//
// The log destination is case insensitive. If the string is not recognized, the default log destination is returned.
func GetLogDestinationFromString(logDestination string) LogDestination {
	switch strings.ToLower(logDestination) {
	case "stdout":
		return LogDestinationStdout
	case "stderr":
		return LogDestinationStderr
	}
	return DefaultLogDestination
}

// GetDefaultLogDestination returns the default log destination.
//
// The default log destination is defined by the environment variable LOG_DESTINATION.
// If the environment variable is not set or empty, the default log destination is stderr.
func GetDefaultLogDestination() LogDestination {
	logDestinationEnvVarMutex.RLock()
	defer logDestinationEnvVarMutex.RUnlock()
	logDestinationAsString := strings.TrimSpace(os.Getenv(logDestinationEnvVar))
	return GetLogDestinationFromString(logDestinationAsString)
}

func (ld LogDestination) getFile() *os.File {
	switch ld {
	case LogDestinationStdout:
		return os.Stdout
	case LogDestinationStderr:
		return os.Stderr
	default:
		panic(fmt.Sprintf("unknown log destination: %s", ld))
	}
}

func getDestination(destination *LogDestination) LogDestination {
	if destination == nil {
		return GetDefaultLogDestination()
	}
	return *destination
}
