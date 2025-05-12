package analyzer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SecurityIssue represents a security issue found in the Docker image
type SecurityIssue struct {
	Type        string `json:"type"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// SecurityResult represents the result of a security analysis
type SecurityResult struct {
	Issues         []SecurityIssue `json:"issues"`
	TotalIssues    int             `json:"totalIssues"`
	HighSeverity   int             `json:"highSeverity"`
	MediumSeverity int             `json:"mediumSeverity"`
	LowSeverity    int             `json:"lowSeverity"`
}

// AnalyzeSecurity analyzes the security of a Docker image
func AnalyzeSecurity(imagePath string) (*SecurityResult, error) {
	result := &SecurityResult{}

	// Check for sensitive files
	sensitiveFiles, err := findSensitiveFiles(imagePath)
	if err != nil {
		return nil, err
	}
	result.Issues = append(result.Issues, sensitiveFiles...)

	// Check for hardcoded secrets
	secrets, err := findHardcodedSecrets(imagePath)
	if err != nil {
		return nil, err
	}
	result.Issues = append(result.Issues, secrets...)

	// Check for vulnerable packages
	// Note: This is a simplified implementation. In a real-world scenario,
	// you would use a vulnerability database like CVE or a service like Trivy.
	vulnPackages, err := findVulnerablePackages(imagePath)
	if err != nil {
		return nil, err
	}
	result.Issues = append(result.Issues, vulnPackages...)

	// Count issues by severity
	for _, issue := range result.Issues {
		switch issue.Severity {
		case "HIGH":
			result.HighSeverity++
		case "MEDIUM":
			result.MediumSeverity++
		case "LOW":
			result.LowSeverity++
		}
	}
	result.TotalIssues = len(result.Issues)

	return result, nil
}

// findSensitiveFiles finds sensitive files in the Docker image
func findSensitiveFiles(imagePath string) ([]SecurityIssue, error) {
	var issues []SecurityIssue

	// Define patterns for sensitive files
	sensitivePatterns := []struct {
		pattern  string
		severity string
		desc     string
	}{
		{".env", "HIGH", "Environment file may contain sensitive information"},
		{".aws/credentials", "HIGH", "AWS credentials file"},
		{".ssh/id_rsa", "HIGH", "SSH private key"},
		{".ssh/id_dsa", "HIGH", "SSH private key"},
		{".ssh/id_ecdsa", "HIGH", "SSH private key"},
		{".ssh/id_ed25519", "HIGH", "SSH private key"},
		{"config.json", "MEDIUM", "Configuration file may contain sensitive information"},
		{"credentials.json", "HIGH", "Credentials file"},
		{"secrets.json", "HIGH", "Secrets file"},
		{"password", "HIGH", "Password file"},
		{".npmrc", "MEDIUM", "NPM configuration file may contain tokens"},
		{".dockercfg", "MEDIUM", "Docker configuration file may contain credentials"},
		{".docker/config.json", "MEDIUM", "Docker configuration file may contain credentials"},
		{"id_rsa", "HIGH", "SSH private key"},
		{"id_dsa", "HIGH", "SSH private key"},
		{"id_ecdsa", "HIGH", "SSH private key"},
		{"id_ed25519", "HIGH", "SSH private key"},
	}

	// Walk the directory tree
	err := filepath.Walk(imagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get the relative path
		relPath, err := filepath.Rel(imagePath, path)
		if err != nil {
			return err
		}

		// Check if the file matches any sensitive pattern
		for _, pattern := range sensitivePatterns {
			if strings.Contains(strings.ToLower(relPath), pattern.pattern) {
				issues = append(issues, SecurityIssue{
					Type:        "SENSITIVE_FILE",
					Path:        relPath,
					Description: pattern.desc,
					Severity:    pattern.severity,
				})
				break
			}
		}

		return nil
	})

	return issues, err
}

// findHardcodedSecrets finds hardcoded secrets in files
func findHardcodedSecrets(imagePath string) ([]SecurityIssue, error) {
	var issues []SecurityIssue

	// Define patterns for hardcoded secrets
	secretPatterns := []struct {
		regex    *regexp.Regexp
		severity string
		desc     string
	}{
		{regexp.MustCompile(`(?i)password\s*=\s*['"]([^'"]{8,})['"]`), "HIGH", "Hardcoded password"},
		{regexp.MustCompile(`(?i)passwd\s*=\s*['"]([^'"]{8,})['"]`), "HIGH", "Hardcoded password"},
		{regexp.MustCompile(`(?i)pwd\s*=\s*['"]([^'"]{8,})['"]`), "HIGH", "Hardcoded password"},
		{regexp.MustCompile(`(?i)secret\s*=\s*['"]([^'"]{8,})['"]`), "HIGH", "Hardcoded secret"},
		{regexp.MustCompile(`(?i)api[_-]?key\s*=\s*['"]([^'"]{8,})['"]`), "HIGH", "Hardcoded API key"},
		{regexp.MustCompile(`(?i)access[_-]?key\s*=\s*['"]([^'"]{8,})['"]`), "HIGH", "Hardcoded access key"},
		{regexp.MustCompile(`(?i)token\s*=\s*['"]([^'"]{8,})['"]`), "HIGH", "Hardcoded token"},
		{regexp.MustCompile(`(?i)aws[_-]?access[_-]?key[_-]?id\s*=\s*['"]([^'"]{16,})['"]`), "HIGH", "Hardcoded AWS access key"},
		{regexp.MustCompile(`(?i)aws[_-]?secret[_-]?access[_-]?key\s*=\s*['"]([^'"]{16,})['"]`), "HIGH", "Hardcoded AWS secret key"},
	}

	// Walk the directory tree
	err := filepath.Walk(imagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and binary files
		if info.IsDir() || isBinaryFile(path) {
			return nil
		}

		// Get the relative path
		relPath, err := filepath.Rel(imagePath, path)
		if err != nil {
			return err
		}

		// Read the file
		file, err := os.Open(path)
		if err != nil {
			return nil // Skip files that can't be opened
		}
		defer file.Close()

		// Scan the file line by line
		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// Check if the line matches any secret pattern
			for _, pattern := range secretPatterns {
				if pattern.regex.MatchString(line) {
					issues = append(issues, SecurityIssue{
						Type:        "HARDCODED_SECRET",
						Path:        fmt.Sprintf("%s:%d", relPath, lineNum),
						Description: pattern.desc,
						Severity:    pattern.severity,
					})
					break
				}
			}
		}

		return scanner.Err()
	})

	return issues, err
}

// findVulnerablePackages finds vulnerable packages in the Docker image
// This is a simplified implementation. In a real-world scenario,
// you would use a vulnerability database like CVE or a service like Trivy.
func findVulnerablePackages(imagePath string) ([]SecurityIssue, error) {
	var issues []SecurityIssue

	// Check for package files
	packageFiles := []struct {
		pattern string
		check   func(string) ([]SecurityIssue, error)
	}{
		{"package.json", checkNodePackages},
		{"requirements.txt", checkPythonPackages},
		{"Gemfile.lock", checkRubyPackages},
		{"go.mod", checkGoPackages},
	}

	// Walk the directory tree
	err := filepath.Walk(imagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get the base name of the file
		baseName := filepath.Base(path)

		// Check if the file is a package file
		for _, pf := range packageFiles {
			if baseName == pf.pattern {
				// Get the relative path
				relPath, err := filepath.Rel(imagePath, path)
				if err != nil {
					return err
				}

				// Check for vulnerable packages
				pkgIssues, err := pf.check(path)
				if err != nil {
					return nil // Skip files that can't be checked
				}

				// Update the path of each issue
				for i := range pkgIssues {
					pkgIssues[i].Path = relPath
				}

				issues = append(issues, pkgIssues...)
			}
		}

		return nil
	})

	return issues, err
}

// checkNodePackages checks for vulnerable Node.js packages
func checkNodePackages(path string) ([]SecurityIssue, error) {
	var issues []SecurityIssue

	// Read the package.json file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse the JSON
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	// Define known vulnerable packages (simplified)
	// In a real-world scenario, you would use a vulnerability database
	vulnPackages := map[string]struct {
		version    string
		severity   string
		desc       string
		fixVersion string
	}{
		"lodash":     {"<4.17.21", "HIGH", "Prototype Pollution in lodash", ">=4.17.21"},
		"minimist":   {"<1.2.6", "HIGH", "Prototype Pollution in minimist", ">=1.2.6"},
		"node-fetch": {"<2.6.7", "HIGH", "Exposure of Sensitive Information in node-fetch", ">=2.6.7"},
		"axios":      {"<0.21.1", "HIGH", "Server-Side Request Forgery in axios", ">=0.21.1"},
	}

	// Check dependencies
	for name, version := range pkg.Dependencies {
		if vuln, ok := vulnPackages[name]; ok {
			// This is a simplified version check
			// In a real-world scenario, you would use semver comparison
			if strings.HasPrefix(version, "^") || strings.HasPrefix(version, "~") {
				version = version[1:]
			}
			if version < vuln.version {
				issues = append(issues, SecurityIssue{
					Type:        "VULNERABLE_PACKAGE",
					Path:        fmt.Sprintf("package.json: %s@%s", name, version),
					Description: fmt.Sprintf("%s. Update to %s", vuln.desc, vuln.fixVersion),
					Severity:    vuln.severity,
				})
			}
		}
	}

	// Check dev dependencies
	for name, version := range pkg.DevDependencies {
		if vuln, ok := vulnPackages[name]; ok {
			// This is a simplified version check
			// In a real-world scenario, you would use semver comparison
			if strings.HasPrefix(version, "^") || strings.HasPrefix(version, "~") {
				version = version[1:]
			}
			if version < vuln.version {
				issues = append(issues, SecurityIssue{
					Type:        "VULNERABLE_PACKAGE",
					Path:        fmt.Sprintf("package.json: %s@%s (dev)", name, version),
					Description: fmt.Sprintf("%s. Update to %s", vuln.desc, vuln.fixVersion),
					Severity:    vuln.severity,
				})
			}
		}
	}

	return issues, nil
}

// checkPythonPackages checks for vulnerable Python packages
func checkPythonPackages(path string) ([]SecurityIssue, error) {
	var issues []SecurityIssue

	// Read the requirements.txt file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Define known vulnerable packages (simplified)
	// In a real-world scenario, you would use a vulnerability database
	vulnPackages := map[string]struct {
		version    string
		severity   string
		desc       string
		fixVersion string
	}{
		"django":   {"<3.2.14", "HIGH", "SQL Injection in Django", ">=3.2.14"},
		"flask":    {"<2.0.1", "MEDIUM", "Open Redirect in Flask", ">=2.0.1"},
		"requests": {"<2.26.0", "MEDIUM", "CRLF Injection in Requests", ">=2.26.0"},
		"pillow":   {"<9.0.0", "HIGH", "Buffer Overflow in Pillow", ">=9.0.0"},
	}

	// Scan the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// Parse the package name and version
		parts := strings.Split(line, "==")
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		version := strings.TrimSpace(parts[1])

		// Check if the package is vulnerable
		if vuln, ok := vulnPackages[name]; ok {
			// This is a simplified version check
			// In a real-world scenario, you would use semver comparison
			if version < vuln.version {
				issues = append(issues, SecurityIssue{
					Type:        "VULNERABLE_PACKAGE",
					Path:        fmt.Sprintf("requirements.txt: %s==%s", name, version),
					Description: fmt.Sprintf("%s. Update to %s", vuln.desc, vuln.fixVersion),
					Severity:    vuln.severity,
				})
			}
		}
	}

	return issues, scanner.Err()
}

// checkRubyPackages checks for vulnerable Ruby packages
func checkRubyPackages(path string) ([]SecurityIssue, error) {
	// Simplified implementation
	// In a real-world scenario, you would parse the Gemfile.lock and check against a vulnerability database
	return []SecurityIssue{}, nil
}

// checkGoPackages checks for vulnerable Go packages
func checkGoPackages(path string) ([]SecurityIssue, error) {
	// Simplified implementation
	// In a real-world scenario, you would parse the go.mod file and check against a vulnerability database
	return []SecurityIssue{}, nil
}

// isBinaryFile checks if a file is binary
func isBinaryFile(path string) bool {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(path))
	binaryExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
		".pdf": true, ".zip": true, ".tar": true, ".gz": true, ".tgz": true,
		".rar": true, ".7z": true, ".exe": true, ".dll": true, ".so": true,
		".dylib": true, ".bin": true, ".dat": true, ".iso": true, ".img": true,
	}
	if binaryExts[ext] {
		return true
	}

	// Check file content (read first few bytes)
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read the first 512 bytes
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil || n == 0 {
		return false
	}

	// Check for null bytes (common in binary files)
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true
		}
	}

	return false
}
