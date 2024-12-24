package tests

import (
	"os"
	"path/filepath"
	"testing"

	"llm_cli/config" // Import the config package
)

func TestGetConfigPath(t *testing.T) {
	path := config.GetConfigPath()  // Need to export this function
	if path == "" {
		t.Error("Expected config path, got empty string")
	}

	if filepath.Base(path) != config.ConfigFile {  // Need to export this constant
		t.Errorf("Expected path to end with %s, got %s", config.ConfigFile, filepath.Base(path))
	}
}

func TestBackupAndRestore(t *testing.T) {
	// Create a temporary test file
	tmpFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testData := []byte(`{"test": "data"}`)
	if err := os.WriteFile(tmpFile.Name(), testData, 0644); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	if err := config.BackupConfig(tmpFile.Name()); err != nil {  // Need to export this function
		t.Errorf("BackupConfig failed: %v", err)
	}

	if err := os.WriteFile(tmpFile.Name(), []byte(`{"modified": "data"}`), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	if err := config.RestoreConfig(tmpFile.Name()); err != nil {  // Need to export this function
		t.Errorf("RestoreConfig failed: %v", err)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read restored file: %v", err)
	}
	if string(content) != string(testData) {
		t.Errorf("Restored content doesn't match. Expected %s, got %s", testData, content)
	}

	os.Remove(tmpFile.Name() + ".backup")
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid json",
			content: `{"key": "value"}`,
			wantErr: false,
		},
		{
			name:    "empty object",
			content: `{}`,
			wantErr: false,
		},
		{
			name:    "invalid json",
			content: `{"key": "value"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "config_test_*.json")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if err := os.WriteFile(tmpFile.Name(), []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write test data: %v", err)
			}

			err = config.ValidateConfig(tmpFile.Name())  // Need to export this function
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} 