package tests

import (
	"os"
	"testing"

	"llm_cli/utils"
)

func TestHistory(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "history_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set HOME environment variable to use temp directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Initialize history
	history, err := utils.NewHistory()
	if err != nil {
		t.Fatalf("Failed to create history: %v", err)
	}
	defer history.Close()

	// Test Push
	t.Run("Push", func(t *testing.T) {
		err := history.Push("test-assistant", "user", "hello")
		if err != nil {
			t.Errorf("Push failed: %v", err)
		}
	})

	// Test Fetch
	t.Run("Fetch", func(t *testing.T) {
		records, err := history.Fetch("test-assistant", 10)
		if err != nil {
			t.Errorf("Fetch failed: %v", err)
		}

		if len(records) != 1 {
			t.Errorf("Expected 1 record, got %d", len(records))
		}

		if records[0].Assistant != "test-assistant" {
			t.Errorf("Expected assistant 'test-assistant', got '%s'", records[0].Assistant)
		}

		if records[0].Role != "user" {
			t.Errorf("Expected role 'user', got '%s'", records[0].Role)
		}

		if records[0].Content != "hello" {
			t.Errorf("Expected content 'hello', got '%s'", records[0].Content)
		}
	})

	// Test multiple records
	t.Run("Multiple Records", func(t *testing.T) {
		// Add more records
		history.Push("test-assistant", "assistant", "hi there")
		history.Push("test-assistant", "user", "how are you")

		records, err := history.Fetch("test-assistant", 3)
		if err != nil {
			t.Errorf("Fetch failed: %v", err)
		}

		if len(records) != 3 {
			t.Errorf("Expected 3 records, got %d", len(records))
		}
	})

	// Test Clear
	t.Run("Clear", func(t *testing.T) {
		err := history.Clear("test-assistant")
		if err != nil {
			t.Errorf("Clear failed: %v", err)
		}

		records, err := history.Fetch("test-assistant", 10)
		if err != nil {
			t.Errorf("Fetch failed: %v", err)
		}

		if len(records) != 0 {
			t.Errorf("Expected 0 records after clear, got %d", len(records))
		}
	})

	// Test limit in Fetch
	t.Run("Fetch Limit", func(t *testing.T) {
		// Add multiple records
		for i := 0; i < 5; i++ {
			history.Push("test-assistant", "user", "message")
		}

		records, err := history.Fetch("test-assistant", 3)
		if err != nil {
			t.Errorf("Fetch failed: %v", err)
		}

		if len(records) != 3 {
			t.Errorf("Expected 3 records with limit, got %d", len(records))
		}
	})
} 