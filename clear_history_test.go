package main

import (
	"os"
	"testing"
)

func TestClearUserHistory(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "iphistory_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Override historyFile variable
	originalHistoryFile := historyFile
	historyFile = tmpFile.Name()
	defer func() { historyFile = originalHistoryFile }()

	// Write initial data
	initialData := "2023-10-25 10:00:00|192.168.1.1\n" +
		"2023-10-25 10:01:00|10.0.0.1\n" +
		"\n" + // Test empty line preservation
		"2023-10-25 10:02:00|192.168.1.2\n"

	if err := os.WriteFile(historyFile, []byte(initialData), 0644); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Run clearUserHistory
	err = clearUserHistory("10.0.0.1")
	if err != nil {
		t.Fatalf("clearUserHistory failed: %v", err)
	}

	// Read result and verify
	resultData, err := os.ReadFile(historyFile)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	expectedData := "2023-10-25 10:00:00|192.168.1.1\n" +
		"\n" + // Empty line should be preserved
		"2023-10-25 10:02:00|192.168.1.2\n"

	if string(resultData) != expectedData {
		t.Errorf("Expected %q, got %q", expectedData, string(resultData))
	}
}
