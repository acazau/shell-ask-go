package spinner

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNewSpinner(t *testing.T) {
	s := New("Testing...")
	if s.message != "Testing..." {
		t.Errorf("Expected message 'Testing...', got '%s'", s.message)
	}
	if len(s.frames) == 0 {
		t.Error("Frames should not be empty")
	}
}

func TestStartAndStop(t *testing.T) {
	s := New("Testing...")
	s.Start()
	time.Sleep(150 * time.Millisecond) // Let the spinner run for a bit
	s.Stop()

	select {
	case <-s.stop:
	default:
		t.Error("Stop channel should be closed after Stop()")
	}
}

func TestNewSpinnerWithBriandowns(t *testing.T) {
	s := NewSpinner("Testing...")
	if s.Suffix != " Testing..." {
		t.Errorf("Expected suffix ' Testing...', got '%s'", s.Suffix)
	}
	if s.Delay != 100*time.Millisecond {
		t.Errorf("Expected delay 100ms, got %v", s.Delay)
	}
}

func TestRunWithSpinner(t *testing.T) {
	var testErr error
	err := RunWithSpinner(&CLI{}, "testProvider", "testPrompt", false, func() error {
		testErr = fmt.Errorf("test error")
		return testErr
	})

	if err == nil {
		t.Error("Expected an error to be returned")
	}

	if err != nil && !strings.Contains(err.Error(), "error during execution: test error") {
		t.Errorf("Expected error message 'error during execution: test error', got '%s'", err.Error())
	}
}
