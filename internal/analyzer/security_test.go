package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "security-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files that should trigger security issues
	sensitiveFiles := []struct {
		path     string
		content  string
		severity string
	}{
		{".env", "DB_PASSWORD=secret123", "HIGH"},
		{".ssh/id_rsa", "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...", "HIGH"},
		{"config.json", "{\"password\": \"secret123\"}", "MEDIUM"},
	}

	// Create files with hardcoded secrets
	secretFiles := []struct {
		path     string
		content  string
		severity string
	}{
		{"config.js", "const password = 'supersecretpassword';", "HIGH"},
		{"api.js", "const apiKey = 'abcdef123456789';", "HIGH"},
		{"aws.js", "const awsAccessKeyId = 'AKIAIOSFODNN7EXAMPLE';", "HIGH"},
	}

	// Create package files with vulnerable packages
	packageFiles := []struct {
		path    string
		content string
	}{
		{"package.json", `{
			"name": "test-app",
			"version": "1.0.0",
			"dependencies": {
				"lodash": "4.17.20",
				"axios": "0.21.0"
			}
		}`},
		{"requirements.txt", "django==3.2.0\nflask==2.0.0\nrequests==2.25.0"},
	}

	// Create all the test files
	for _, file := range sensitiveFiles {
		dirPath := filepath.Dir(filepath.Join(tempDir, file.path))
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %q: %v", dirPath, err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, file.path), []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", file.path, err)
		}
	}

	for _, file := range secretFiles {
		dirPath := filepath.Dir(filepath.Join(tempDir, file.path))
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %q: %v", dirPath, err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, file.path), []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", file.path, err)
		}
	}

	for _, file := range packageFiles {
		dirPath := filepath.Dir(filepath.Join(tempDir, file.path))
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %q: %v", dirPath, err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, file.path), []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", file.path, err)
		}
	}

	// Run the security analysis
	result, err := AnalyzeSecurity(tempDir)
	if err != nil {
		t.Fatalf("AnalyzeSecurity failed: %v", err)
	}

	// Check that we have security issues
	if result.TotalIssues == 0 {
		t.Errorf("Expected security issues, but found none")
	}

	// Check that we have the expected number of issues
	expectedIssueCount := len(sensitiveFiles) + len(secretFiles) + 2 // 2 vulnerable packages
	if result.TotalIssues < expectedIssueCount {
		t.Errorf("Expected at least %d security issues, but found %d", expectedIssueCount, result.TotalIssues)
	}

	// Check that we have high severity issues
	if result.HighSeverity == 0 {
		t.Errorf("Expected high severity issues, but found none")
	}

	// Check that we have medium severity issues
	if result.MediumSeverity == 0 {
		t.Errorf("Expected medium severity issues, but found none")
	}

	// Check that we have the expected issue types
	issueTypes := make(map[string]bool)
	for _, issue := range result.Issues {
		issueTypes[issue.Type] = true
	}

	expectedTypes := []string{"SENSITIVE_FILE", "HARDCODED_SECRET", "VULNERABLE_PACKAGE"}
	for _, expectedType := range expectedTypes {
		if !issueTypes[expectedType] {
			t.Errorf("Expected issue type %q, but it was not found", expectedType)
		}
	}
}

func TestFindSensitiveFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "sensitive-files-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files that should trigger security issues
	sensitiveFiles := []struct {
		path     string
		content  string
		severity string
	}{
		{".env", "DB_PASSWORD=secret123", "HIGH"},
		{".ssh/id_rsa", "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...", "HIGH"},
		{"config.json", "{\"password\": \"secret123\"}", "MEDIUM"},
		{".npmrc", "//registry.npmjs.org/:_authToken=npm_token", "MEDIUM"},
		{".docker/config.json", "{\"auths\": {\"registry\": {\"auth\": \"base64token\"}}}", "MEDIUM"},
	}

	// Create all the test files
	for _, file := range sensitiveFiles {
		dirPath := filepath.Dir(filepath.Join(tempDir, file.path))
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %q: %v", dirPath, err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, file.path), []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", file.path, err)
		}
	}

	// Run the sensitive files check
	issues, err := findSensitiveFiles(tempDir)
	if err != nil {
		t.Fatalf("findSensitiveFiles failed: %v", err)
	}

	// Check that we found all the sensitive files
	if len(issues) != len(sensitiveFiles) {
		t.Errorf("Expected %d sensitive files, but found %d", len(sensitiveFiles), len(issues))
	}

	// Check that all issues have the correct type
	for _, issue := range issues {
		if issue.Type != "SENSITIVE_FILE" {
			t.Errorf("Expected issue type SENSITIVE_FILE, but got %q", issue.Type)
		}
	}

	// Check that we have the expected severity levels
	severityCounts := make(map[string]int)
	for _, issue := range issues {
		severityCounts[issue.Severity]++
	}

	if severityCounts["HIGH"] < 2 {
		t.Errorf("Expected at least 2 HIGH severity issues, but found %d", severityCounts["HIGH"])
	}

	if severityCounts["MEDIUM"] < 3 {
		t.Errorf("Expected at least 3 MEDIUM severity issues, but found %d", severityCounts["MEDIUM"])
	}
}

func TestFindHardcodedSecrets(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "hardcoded-secrets-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create files with hardcoded secrets
	secretFiles := []struct {
		path    string
		content string
	}{
		{"config.js", "const password = 'supersecretpassword';"},
		{"api.js", "const apiKey = 'abcdef123456789';"},
		{"aws.js", "const awsAccessKeyId = 'AKIAIOSFODNN7EXAMPLE';"},
		{"aws2.js", "const awsSecretAccessKey = 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY';"},
		{"token.js", "const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9';"},
	}

	// Create all the test files
	for _, file := range secretFiles {
		dirPath := filepath.Dir(filepath.Join(tempDir, file.path))
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %q: %v", dirPath, err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, file.path), []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", file.path, err)
		}
	}

	// Run the hardcoded secrets check
	issues, err := findHardcodedSecrets(tempDir)
	if err != nil {
		t.Fatalf("findHardcodedSecrets failed: %v", err)
	}

	// Check that we found all the hardcoded secrets
	if len(issues) != len(secretFiles) {
		t.Errorf("Expected %d hardcoded secrets, but found %d", len(secretFiles), len(issues))
	}

	// Check that all issues have the correct type
	for _, issue := range issues {
		if issue.Type != "HARDCODED_SECRET" {
			t.Errorf("Expected issue type HARDCODED_SECRET, but got %q", issue.Type)
		}
	}

	// Check that all issues have HIGH severity
	for _, issue := range issues {
		if issue.Severity != "HIGH" {
			t.Errorf("Expected severity HIGH, but got %q", issue.Severity)
		}
	}
}

func TestFindVulnerablePackages(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "vulnerable-packages-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create package files with vulnerable packages
	packageFiles := []struct {
		path    string
		content string
	}{
		{"package.json", `{
			"name": "test-app",
			"version": "1.0.0",
			"dependencies": {
				"lodash": "4.17.20",
				"axios": "0.21.0"
			},
			"devDependencies": {
				"minimist": "1.2.5"
			}
		}`},
		{"requirements.txt", "django==3.2.0\nflask==2.0.0\nrequests==2.25.0\npillow==8.0.0"},
	}

	// Create all the test files
	for _, file := range packageFiles {
		dirPath := filepath.Dir(filepath.Join(tempDir, file.path))
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %q: %v", dirPath, err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, file.path), []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", file.path, err)
		}
	}

	// Run the vulnerable packages check
	issues, err := findVulnerablePackages(tempDir)
	if err != nil {
		t.Fatalf("findVulnerablePackages failed: %v", err)
	}

	// Check that we found vulnerable packages
	if len(issues) == 0 {
		t.Errorf("Expected vulnerable packages, but found none")
	}

	// Check that all issues have the correct type
	for _, issue := range issues {
		if issue.Type != "VULNERABLE_PACKAGE" {
			t.Errorf("Expected issue type VULNERABLE_PACKAGE, but got %q", issue.Type)
		}
	}

	// Check that we found specific vulnerable packages
	vulnerablePackages := []string{"lodash", "axios", "django", "pillow"}
	foundPackages := make(map[string]bool)

	for _, issue := range issues {
		for _, pkg := range vulnerablePackages {
			if issue.Path != "" && issue.Description != "" && issue.Severity != "" {
				foundPackages[pkg] = true
			}
		}
	}

	for _, pkg := range vulnerablePackages {
		if !foundPackages[pkg] {
			t.Errorf("Expected to find vulnerable package %q, but it was not detected", pkg)
		}
	}
}
