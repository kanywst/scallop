package docker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsDockerImageName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"nginx", true},
		{"nginx:latest", true},
		{"docker.io/library/nginx", true},
		{"docker.io/library/nginx:latest", true},
		{"registry.example.com/myapp:v1.0.0", true},
		{"user/repo:tag", true},
		{"/path/to/image.tar", false},
		{"./image.tar", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsDockerImageName(test.name)
		if result != test.expected {
			t.Errorf("IsDockerImageName(%q) = %v, expected %v", test.name, result, test.expected)
		}
	}
}

func TestExtractImage(t *testing.T) {
	// This is a simplified test that doesn't actually extract a Docker image
	// In a real test, you would use a mock Docker client or a small test image

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "docker-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock image tar file
	mockImagePath := filepath.Join(tempDir, "mock-image.tar")
	if err := os.WriteFile(mockImagePath, []byte("mock tar content"), 0644); err != nil {
		t.Fatalf("Failed to create mock image file: %v", err)
	}

	// Create a destination directory
	destDir := filepath.Join(tempDir, "dest")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}

	// Test with a non-existent image path
	_, err = ExtractImage("non-existent-image.tar", destDir)
	if err == nil {
		t.Errorf("ExtractImage with non-existent image should fail")
	}

	// Note: We can't fully test ExtractImage without a real Docker image
	// or mocking the Docker client, so we'll skip the actual extraction test
}

// TestExtractFromTarFile is a simplified test for the extractFromTarFile function
func TestExtractFromTarFile(t *testing.T) {
	// This is a simplified test that doesn't actually extract a Docker image
	// In a real test, you would use a mock Docker client or a small test image

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "extract-tar-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock tar file
	mockTarPath := filepath.Join(tempDir, "mock-image.tar")
	if err := os.WriteFile(mockTarPath, []byte("mock tar content"), 0644); err != nil {
		t.Fatalf("Failed to create mock tar file: %v", err)
	}

	// Create a destination directory
	destDir := filepath.Join(tempDir, "dest")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}

	// Test with a non-existent tar file
	_, err = extractFromTarFile("non-existent-file.tar", destDir)
	if err == nil {
		t.Errorf("extractFromTarFile with non-existent file should fail")
	}

	// Note: We can't fully test extractFromTarFile without a real Docker image
	// or mocking the tar extraction, so we'll skip the actual extraction test
}
