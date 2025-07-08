package filter_test

import (
	"docklog/internal/filter"
	"testing"
)

func TestIsErrorLog(t *testing.T) {
	testCases := []struct {
		name      string
		log       string
		logStream string
		expected  bool
	}{
		{"stdout normal log", "Application started", "stdout", false},
		{"stdout with error keyword", "Error occurred", "stdout", true},
		{"stderr always error", "Normal log", "stderr", true},
		{"stdout with fail keyword", "Operation failed", "stdout", true},
		{"stdout with exception", "Exception in thread", "stdout", true},
		{"stdout with fatal", "Fatal error", "stdout", true},
		{"case insensitive", "ERROR: something wrong", "stdout", true},
		{"empty log stdout", "", "stdout", false},
		{"empty log stderr", "", "stderr", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.IsErrorLog(tc.log, tc.logStream)
			if result != tc.expected {
				t.Errorf("IsErrorLog(%q, %q) = %v; expected %v",
					tc.log, tc.logStream, result, tc.expected)
			}
		})
	}
}

func TestIsErrorLog_StderrAlwaysError(t *testing.T) {
	result := filter.IsErrorLog("normal application log", "stderr")
	if !result {
		t.Error("stderr logs should always be considered errors")
	}
}

func TestIsErrorLog_ErrorKeywords(t *testing.T) {
	errorKeywords := []string{"error", "fail", "exception", "fatal", "failed", "problem"}

	for _, keyword := range errorKeywords {
		log := "Something " + keyword + " happened"
		result := filter.IsErrorLog(log, "stdout")
		if !result {
			t.Errorf("Should detect %q as error keyword in log: %q", keyword, log)
		}
	}
}

func TestIsErrorLog_CaseInsensitive(t *testing.T) {
	testCases := []string{"ERROR", "Error", "error", "eRrOr"}

	for _, errorWord := range testCases {
		log := "Something " + errorWord + " happened"
		result := filter.IsErrorLog(log, "stdout")
		if !result {
			t.Errorf("Should detect error regardless of case: %q", log)
		}
	}
}

func TestIsErrorLog_NoFalsePositives(t *testing.T) {
	normalLogs := []string{
		"Application started successfully",
		"Processing 1000 records",
		"User login successful",
		"Database connection established",
		"Server listening on port 8080",
	}

	for _, log := range normalLogs {
		result := filter.IsErrorLog(log, "stdout")
		if result {
			t.Errorf("Should not detect error in normal log: %q", log)
		}
	}
}

func TestIsErrorLog_EdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		log       string
		logStream string
		expected  bool
	}{
		{"empty log stdout", "", "stdout", false},
		{"empty log stderr", "", "stderr", true},
		{"whitespace only stdout", "   ", "stdout", false},
		{"whitespace only stderr", "   ", "stderr", true},
		{"error as substring", "terrific", "stdout", false},
		{"error with punctuation", "Error!", "stdout", true},
		{"multiple error words", "fatal error exception", "stdout", true},
		{"unknown stream type", "error occurred", "unknown", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.IsErrorLog(tc.log, tc.logStream)
			if result != tc.expected {
				t.Errorf("IsErrorLog(%q, %q) = %v; expected %v",
					tc.log, tc.logStream, result, tc.expected)
			}
		})
	}
}
