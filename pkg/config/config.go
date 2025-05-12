package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the configuration for the scallop tool
type Config struct {
	// Default output format (text or json)
	DefaultOutputFormat string `json:"defaultOutputFormat"`

	// Default verbosity level
	Verbose bool `json:"verbose"`

	// Security scan settings
	Security SecurityConfig `json:"security"`

	// Size analysis settings
	Size SizeConfig `json:"size"`
}

// SecurityConfig represents the security scan configuration
type SecurityConfig struct {
	// Enable security scanning
	Enabled bool `json:"enabled"`

	// Minimum severity level to report (LOW, MEDIUM, HIGH)
	MinSeverity string `json:"minSeverity"`

	// Custom patterns for sensitive files
	SensitivePatterns []string `json:"sensitivePatterns"`

	// Custom patterns for hardcoded secrets
	SecretPatterns []string `json:"secretPatterns"`
}

// SizeConfig represents the size analysis configuration
type SizeConfig struct {
	// Enable size analysis
	Enabled bool `json:"enabled"`

	// Number of largest files to report
	TopFilesCount int `json:"topFilesCount"`

	// Number of largest directories to report
	TopDirsCount int `json:"topDirsCount"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultOutputFormat: "text",
		Verbose:             false,
		Security: SecurityConfig{
			Enabled:     true,
			MinSeverity: "LOW",
		},
		Size: SizeConfig{
			Enabled:       true,
			TopFilesCount: 10,
			TopDirsCount:  5,
		},
	}
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	// Use default config if no path is provided
	if path == "" {
		return DefaultConfig(), nil
	}

	// Read the config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the config file
	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, path string) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Marshal the config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write the config file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	// Check if the config file exists in the current directory
	if _, err := os.Stat("scallop.json"); err == nil {
		return "scallop.json"
	}

	// Check if the config file exists in the user's home directory
	home, err := os.UserHomeDir()
	if err == nil {
		path := filepath.Join(home, ".config", "scallop", "config.json")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Return the default path
	return ""
}
