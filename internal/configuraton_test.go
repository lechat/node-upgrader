package internal

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestReadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Define test data
	testConfig := `
log_level: "info"
account_batch_size: 10
`
	// Write test data to a temporary config file
	configPath := filepath.Join(tmpDir, "config.yaml")
	err := ioutil.WriteFile(configPath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Test reading the config file
	config, err := ReadConfig(tmpDir)
	if err != nil {
		t.Fatalf("ReadConfig returned error: %v", err)
	}

	// Verify the values in the config struct
	expectedConfig := &Config{
		LogLevel:         "info",
		AccountBatchSize: 10,
	}
	if config.LogLevel != expectedConfig.LogLevel {
		t.Errorf("got LogLevel %s, want %s", config.LogLevel, expectedConfig.LogLevel)
	}
	if config.AccountBatchSize != expectedConfig.AccountBatchSize {
		t.Errorf("got AccountBatchSize %d, want %d", config.AccountBatchSize, expectedConfig.AccountBatchSize)
	}
}

func TestReadConfig_Error(t *testing.T) {
	// Test reading a non-existent config file
	_, err := ReadConfig("nonexistent")
	if err == nil {
		t.Error("expected error, got nil")
	}
}
