package filter

import (
	"strings"
)

// returns true if the log is an error
func IsError(log string, logStream string) bool {
	lowerLog := strings.ToLower(log)

	if logStream == "stderr" || strings.Contains(lowerLog, "error") || strings.Contains(lowerLog, "fail") || strings.Contains(lowerLog, "exception") || strings.Contains(lowerLog, "fatal") {
		return true
	}

	return false
}
