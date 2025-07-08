package filter

import (
	"strings"

	"github.com/docker/docker/api/types/events"
)

// returns true if the log is an error
func IsErrorLog(logStr string, logStream string) bool {
	lowerLog := strings.ToLower(logStr)

	if logStream == "stderr" || strings.Contains(lowerLog, "error") || strings.Contains(lowerLog, "fail") || strings.Contains(lowerLog, "exception") || strings.Contains(lowerLog, "fatal") || strings.Contains(lowerLog, "failed") || strings.Contains(lowerLog, "problem") {
		return true
	}

	return false
}

// returns true if the event is a container start or restart
func IsContainerAction(eventType events.Type, action events.Action) bool {
	return eventType == "container" && (action == "start" || action == "restart")
}
