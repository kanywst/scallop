package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	// Get the default configuration
	cfg := DefaultConfig()

	// Check that the default values are set correctly
	if cfg.DefaultOutputFormat != "text" {
		t.Errorf("DefaultOutputFormat = %q, expected %q", cfg.DefaultOutputFormat, "text")
	}

	if cfg.Verbose {
		t.Errorf("Verbose = %v, expected %v", cfg.Verbose, false)
	}

	if !cfg.Security.Enabled {
		t.Errorf("Security.Enabled = %v, expected %v", cfg.Security.Enabled, true)
	}

	if cfg.Security.MinSeverity != "LOW" {
		t.Errorf("Security.MinSeverity = %q, expected %q", cfg.Security.MinSeverity, "LOW")
	}

	if !cfg.Size.Enabled {
		t.Errorf("Size.Enabled = %v, expected %v", cfg.Size.Enabled, true)
	}

	if cfg.Size.TopFilesCount != 10 {
		t.Errorf("Size.TopFilesCount = %d, expected %d", cfg.Size.TopFilesCount, 10)
	}

	if cfg.Size.TopDirsCount != 5 {
		t.Errorf("Size.TopDirsCount = %d, expected %d", cfg.Size.TopDirsCount, 5)
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test configuration file
	configPath := filepath.Join(tempDir, "config.json")
	testConfig := &Config{
		DefaultOutputFormat: "json",
		Verbose:             true,
		Security: SecurityConfig{
			Enabled:     false,
			MinSeverity: "HIGH",
		},
		Size: SizeConfig{
			Enabled:       true,
			TopFilesCount: 20,
			TopDirsCount:  10,
		},
	}

	// Marshal the test configuration to JSON
	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test configuration: %v", err)
	}

	// Write the test configuration to a file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test configuration: %v", err)
	}

	// Load the configuration from the file
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Check that the loaded configuration matches the test configuration
	if loadedConfig.DefaultOutputFormat != testConfig.DefaultOutputFormat {
		t.Errorf("DefaultOutputFormat = %q, expected %q", loadedConfig.DefaultOutputFormat, testConfig.DefaultOutputFormat)
	}

	if loadedConfig.Verbose != testConfig.Verbose {
		t.Errorf("Verbose = %v, expected %v", loadedConfig.Verbose, testConfig.Verbose)
	}

	if loadedConfig.Security.Enabled != testConfig.Security.Enabled {
		t.Errorf("Security.Enabled = %v, expected %v", loadedConfig.Security.Enabled, testConfig.Security.Enabled)
	}

	if loadedConfig.Security.MinSeverity != testConfig.Security.MinSeverity {
		t.Errorf("Security.MinSeverity = %q, expected %q", loadedConfig.Security.MinSeverity, testConfig.Security.MinSeverity)
	}

	if loadedConfig.Size.Enabled != testConfig.Size.Enabled {
		t.Errorf("Size.Enabled = %v, expected %v", loadedConfig.Size.Enabled, testConfig.Size.Enabled)
	}

	if loadedConfig.Size.TopFilesCount != testConfig.Size.TopFilesCount {
		t.Errorf("Size.TopFilesCount = %d, expected %d", loadedConfig.Size.TopFilesCount, testConfig.Size.TopFilesCount)
	}

	if loadedConfig.Size.TopDirsCount != testConfig.Size.TopDirsCount {
		t.Errorf("Size.TopDirsCount = %d, expected %d", loadedConfig.Size.TopDirsCount, testConfig.Size.TopDirsCount)
	}

	// Test loading with an empty path (should return default config)
	defaultConfig, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig with empty path failed: %v", err)
	}

	// Check that the default configuration is returned
	if defaultConfig.DefaultOutputFormat != "text" {
		t.Errorf("DefaultOutputFormat = %q, expected %q", defaultConfig.DefaultOutputFormat, "text")
	}

	// Test loading with a non-existent path (should return an error)
	_, err = LoadConfig("non-existent-file.json")
	if err == nil {
		t.Errorf("LoadConfig with non-existent path should fail")
	}

	// Test loading with an invalid JSON file
	invalidPath := filepath.Join(tempDir, "invalid.json")
	if err := os.WriteFile(invalidPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write invalid JSON file: %v", err)
	}

	_, err = LoadConfig(invalidPath)
	if err == nil {
		t.Errorf("LoadConfig with invalid JSON file should fail")
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test configuration
	testConfig := &Config{
		DefaultOutputFormat: "json",
		Verbose:             true,
		Security: SecurityConfig{
			Enabled:     false,
			MinSeverity: "HIGH",
		},
		Size: SizeConfig{
			Enabled:       true,
			TopFilesCount: 20,
			TopDirsCount:  10,
		},
	}

	// Save the configuration to a file
	configPath := filepath.Join(tempDir, "config.json")
	err = SaveConfig(testConfig, configPath)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Check that the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
	}

	// Load the configuration from the file
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved configuration: %v", err)
	}

	// Check that the loaded configuration matches the test configuration
	if loadedConfig.DefaultOutputFormat != testConfig.DefaultOutputFormat {
		t.Errorf("DefaultOutputFormat = %q, expected %q", loadedConfig.DefaultOutputFormat, testConfig.DefaultOutputFormat)
	}

	if loadedConfig.Verbose != testConfig.Verbose {
		t.Errorf("Verbose = %v, expected %v", loadedConfig.Verbose, testConfig.Verbose)
	}

	// Test saving to a nested directory that doesn't exist
	nestedPath := filepath.Join(tempDir, "nested", "dir", "config.json")
	err = SaveConfig(testConfig, nestedPath)
	if err != nil {
		t.Fatalf("SaveConfig to nested directory failed: %v", err)
	}

	// Check that the file exists
	if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created in nested directory")
	}
}

func TestGetConfigPath(t *testing.T) {
	// Save the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temporary directory: %v", err)
	}
	defer os.Chdir(cwd) // Restore the original working directory

	// Test with no config file (should return empty string)
	path := GetConfigPath()
	if path != "" {
		t.Errorf("GetConfigPath with no config file = %q, expected empty string", path)
	}

	// Create a config file in the current directory
	if err := os.WriteFile("scallop.json", []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Test with a config file in the current directory
	path = GetConfigPath()
	if path != "scallop.json" {
		t.Errorf("GetConfigPath with config file in current directory = %q, expected %q", path, "scallop.json")
	}

	// Remove the config file
	if err := os.Remove("scallop.json"); err != nil {
		t.Fatalf("Failed to remove config file: %v", err)
	}

	// Create a config file in the user's home directory
	// Note: This is a simplified test that doesn't actually create a file in the user's home directory
	// In a real test, you would mock the os.UserHomeDir function or create a temporary home directory
}
