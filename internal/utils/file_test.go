package utils

import (
	"archive/tar"
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "file-exists-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Test with an existing file
	if !FileExists(tempFile.Name()) {
		t.Errorf("FileExists(%q) = false, expected true", tempFile.Name())
	}

	// Test with a non-existent file
	if FileExists("non-existent-file") {
		t.Errorf("FileExists(%q) = true, expected false", "non-existent-file")
	}
}

func TestDirExists(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "dir-exists-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with an existing directory
	if !DirExists(tempDir) {
		t.Errorf("DirExists(%q) = false, expected true", tempDir)
	}

	// Test with a non-existent directory
	if DirExists("non-existent-dir") {
		t.Errorf("DirExists(%q) = true, expected false", "non-existent-dir")
	}

	// Test with a file (should return false)
	tempFile, err := os.CreateTemp("", "dir-exists-test-file-")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if DirExists(tempFile.Name()) {
		t.Errorf("DirExists(%q) = true, expected false for a file", tempFile.Name())
	}
}

func TestGetFileSize(t *testing.T) {
	// Create a temporary file with known content
	tempFile, err := os.CreateTemp("", "file-size-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some content to the file
	content := []byte("test content")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tempFile.Close()

	// Test with an existing file
	size, err := GetFileSize(tempFile.Name())
	if err != nil {
		t.Errorf("GetFileSize(%q) failed: %v", tempFile.Name(), err)
	}
	if size != int64(len(content)) {
		t.Errorf("GetFileSize(%q) = %d, expected %d", tempFile.Name(), size, len(content))
	}

	// Test with a non-existent file
	_, err = GetFileSize("non-existent-file")
	if err == nil {
		t.Errorf("GetFileSize(%q) should fail with a non-existent file", "non-existent-file")
	}
}

func TestGetDirSize(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "dir-size-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some files in the directory
	files := []struct {
		name    string
		content string
	}{
		{"file1.txt", "test content 1"},
		{"file2.txt", "test content 2"},
		{"file3.txt", "test content 3"},
	}

	var expectedSize int64
	for _, file := range files {
		path := filepath.Join(tempDir, file.name)
		if err := os.WriteFile(path, []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", path, err)
		}
		expectedSize += int64(len(file.content))
	}

	// Test with an existing directory
	size, err := GetDirSize(tempDir)
	if err != nil {
		t.Errorf("GetDirSize(%q) failed: %v", tempDir, err)
	}
	if size != expectedSize {
		t.Errorf("GetDirSize(%q) = %d, expected %d", tempDir, size, expectedSize)
	}

	// Test with a non-existent directory
	_, err = GetDirSize("non-existent-dir")
	if err == nil {
		t.Errorf("GetDirSize(%q) should fail with a non-existent directory", "non-existent-dir")
	}
}

func TestCopyFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "copy-file-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a source file
	srcPath := filepath.Join(tempDir, "source.txt")
	content := []byte("test content")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy the file
	dstPath := filepath.Join(tempDir, "destination.txt")
	if err := CopyFile(srcPath, dstPath); err != nil {
		t.Errorf("CopyFile(%q, %q) failed: %v", srcPath, dstPath, err)
	}

	// Check if the destination file exists and has the correct content
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Errorf("Failed to read destination file: %v", err)
	}
	if string(dstContent) != string(content) {
		t.Errorf("Destination file content = %q, expected %q", dstContent, content)
	}

	// Test with a non-existent source file
	err = CopyFile("non-existent-file", dstPath)
	if err == nil {
		t.Errorf("CopyFile(%q, %q) should fail with a non-existent source file", "non-existent-file", dstPath)
	}
}

func TestExtractTar(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "extract-tar-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a tar file
	tarPath := filepath.Join(tempDir, "test.tar")
	tarFile, err := os.Create(tarPath)
	if err != nil {
		t.Fatalf("Failed to create tar file: %v", err)
	}

	// Create a tar writer
	tarWriter := tar.NewWriter(tarFile)

	// Add some files to the tar
	files := []struct {
		name    string
		content string
	}{
		{"file1.txt", "test content 1"},
		{"file2.txt", "test content 2"},
		{"dir/file3.txt", "test content 3"},
	}

	for _, file := range files {
		// Create a tar header
		header := &tar.Header{
			Name: file.name,
			Mode: 0644,
			Size: int64(len(file.content)),
		}

		// Write the header
		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatalf("Failed to write tar header: %v", err)
		}

		// Write the file content
		if _, err := tarWriter.Write([]byte(file.content)); err != nil {
			t.Fatalf("Failed to write tar content: %v", err)
		}
	}

	// Close the tar writer and file
	if err := tarWriter.Close(); err != nil {
		t.Fatalf("Failed to close tar writer: %v", err)
	}
	if err := tarFile.Close(); err != nil {
		t.Fatalf("Failed to close tar file: %v", err)
	}

	// Create a destination directory
	destDir := filepath.Join(tempDir, "extracted")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}

	// Extract the tar file
	if err := ExtractTar(tarPath, destDir); err != nil {
		t.Errorf("ExtractTar(%q, %q) failed: %v", tarPath, destDir, err)
	}

	// Check if the files were extracted correctly
	for _, file := range files {
		path := filepath.Join(destDir, file.name)
		if !FileExists(path) {
			t.Errorf("File %q was not extracted", path)
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read extracted file %q: %v", path, err)
			continue
		}

		if string(content) != file.content {
			t.Errorf("Extracted file %q content = %q, expected %q", path, content, file.content)
		}
	}

	// Test with an invalid tar file
	invalidTarPath := filepath.Join(tempDir, "invalid.tar")
	if err := os.WriteFile(invalidTarPath, []byte("not a valid tar file"), 0644); err != nil {
		t.Fatalf("Failed to create invalid tar file: %v", err)
	}

	err = ExtractTar(invalidTarPath, destDir)
	if err == nil {
		t.Errorf("ExtractTar(%q, %q) should fail with an invalid tar file", invalidTarPath, destDir)
	}

	// Test with a non-existent tar file
	err = ExtractTar("non-existent-file.tar", destDir)
	if err == nil {
		t.Errorf("ExtractTar(%q, %q) should fail with a non-existent tar file", "non-existent-file.tar", destDir)
	}
}
