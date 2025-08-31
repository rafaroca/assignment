package main

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDash0LogsProcessor_TwoLogEntriesOutput(t *testing.T) {
	// Capture stdout to verify output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logIntake := make(chan string, 10)
	processor := &dash0LogsProcessor{
		logStats:       make(map[string]uint64),
		logIntake:      logIntake,
		durationWindow: 50 * time.Millisecond,
	}

	// Start processor in background
	go processor.StartLogProcessing()

	// Add two log entries
	logIntake <- "error-log"
	logIntake <- "info-log"

	// Give processor time to process logs
	time.Sleep(20 * time.Millisecond)

	// Wait for ticker to fire and print stats
	time.Sleep(60 * time.Millisecond)

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout
	
	outputBytes, _ := io.ReadAll(r)
	output := string(outputBytes)

	// Verify output contains expected content
	if !strings.Contains(output, "Log stats:") {
		t.Error("Expected 'Log stats:' in output")
	}
	if !strings.Contains(output, "error-log - 1") {
		t.Error("Expected 'error-log - 1' in output")
	}
	if !strings.Contains(output, "info-log - 1") {
		t.Error("Expected 'info-log - 1' in output")
	}
}
