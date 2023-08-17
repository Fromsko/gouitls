package logs

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestSpiderErrorMethods(t *testing.T) {
	// Initialize SpiderError with an empty message
	err := &SpiderError{msg: ""}

	// Redirect color output for testing
	var output bytes.Buffer
	color.Output = &output

	// Test InfoMsg
	infoMsg := "Info message"
	err.InfoMsg(infoMsg)
	if !strings.Contains(output.String(), infoMsg) {
		t.Errorf("Info message not found in output")
	}

	// Test WarnMsg
	warnMsg := "Warning message"
	err.WarnMsg(warnMsg)
	if !strings.Contains(output.String(), warnMsg) {
		t.Errorf("Warning message not found in output")
	}

	// Test DebugMsg
	debugMsg := "Debug message"
	err.DebugMsg(debugMsg)
	if !strings.Contains(output.String(), debugMsg) {
		t.Errorf("Debug message not found in output")
	}

	// Test ErrorMsg
	errorMsg := "Error message"
	err.ErrorMsg(errorMsg)
	if !strings.Contains(output.String(), errorMsg) {
		t.Errorf("Error message not found in output")
	}
}

func TestChoiceValue(t *testing.T) {
	// Test when a is an empty string
	emptyResult := choiceValue("", []string{"one", "two", "three"})
	if emptyResult != "one - two - three" {
		t.Errorf("Expected 'one - two - three', got '%s'", emptyResult)
	}

	// Test when a is not an empty string
	nonEmptyResult := choiceValue("existing", []string{"one", "two", "three"})
	if nonEmptyResult != "existing" {
		t.Errorf("Expected 'existing', got '%s'", nonEmptyResult)
	}
}
