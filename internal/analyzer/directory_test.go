package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "directory-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a directory structure for testing
	dirs := []string{
		"dir1",
		"dir2",
		"dir1/subdir1",
		"dir2/subdir2",
	}

	files := []struct {
		path    string
		content string
	}{
		{"file1.txt", "test content 1"},
		{"file2.go", "test content 2"},
		{"dir1/file3.txt", "test content 3"},
		{"dir1/subdir1/file4.go", "test content 4"},
		{"dir2/file5.txt", "test content 5"},
		{"dir2/subdir2/file6.go", "test content 6"},
		{"noext", "file without extension"},
	}

	// Create directories
	for _, dir := range dirs {
		dirPath := filepath.Join(tempDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %q: %v", dirPath, err)
		}
	}

	// Create files
	for _, file := range files {
		filePath := filepath.Join(tempDir, file.path)
		if err := os.WriteFile(filePath, []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", filePath, err)
		}
	}

	// Test with verbose=false
	info, err := AnalyzeDirectory(tempDir, false)
	if err != nil {
		t.Fatalf("AnalyzeDirectory failed: %v", err)
	}

	// Check the results
	if info.Path != tempDir {
		t.Errorf("Path = %q, expected %q", info.Path, tempDir)
	}

	expectedFileCount := len(files)
	if info.FileCount != expectedFileCount {
		t.Errorf("FileCount = %d, expected %d", info.FileCount, expectedFileCount)
	}

	expectedDirCount := len(dirs)
	if info.DirCount != expectedDirCount {
		t.Errorf("DirCount = %d, expected %d", info.DirCount, expectedDirCount)
	}

	// Check file types
	expectedFileTypes := map[string]int{
		".txt":           3,
		".go":            3,
		"[no extension]": 1,
	}

	for ext, count := range expectedFileTypes {
		if info.FileTypes[ext] != count {
			t.Errorf("FileTypes[%q] = %d, expected %d", ext, info.FileTypes[ext], count)
		}
	}

	// Files and Dirs should be nil when verbose=false
	if info.Files != nil {
		t.Errorf("Files should be nil when verbose=false")
	}
	if info.Dirs != nil {
		t.Errorf("Dirs should be nil when verbose=false")
	}

	// Test with verbose=true
	verboseInfo, err := AnalyzeDirectory(tempDir, true)
	if err != nil {
		t.Fatalf("AnalyzeDirectory with verbose=true failed: %v", err)
	}

	// Files and Dirs should not be nil when verbose=true
	if verboseInfo.Files == nil {
		t.Errorf("Files should not be nil when verbose=true")
	}
	if verboseInfo.Dirs == nil {
		t.Errorf("Dirs should not be nil when verbose=true")
	}

	// Check that all files are included
	if len(verboseInfo.Files) != expectedFileCount {
		t.Errorf("len(Files) = %d, expected %d", len(verboseInfo.Files), expectedFileCount)
	}

	// Check that all directories are included
	if len(verboseInfo.Dirs) != expectedDirCount {
		t.Errorf("len(Dirs) = %d, expected %d", len(verboseInfo.Dirs), expectedDirCount)
	}
}

func TestGetTopDirectories(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "top-dirs-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a directory structure for testing
	dirs := []string{
		"dir1",
		"dir2",
		"dir3",
		"dir4",
		"dir5",
	}

	files := []struct {
		path    string
		content string
		size    int
	}{
		{"dir1/file1.txt", "a", 1},
		{"dir2/file2.txt", "ab", 2},
		{"dir3/file3.txt", "abc", 3},
		{"dir4/file4.txt", "abcd", 4},
		{"dir5/file5.txt", "abcde", 5},
	}

	// Create directories
	for _, dir := range dirs {
		dirPath := filepath.Join(tempDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %q: %v", dirPath, err)
		}
	}

	// Create files
	for _, file := range files {
		filePath := filepath.Join(tempDir, file.path)
		if err := os.WriteFile(filePath, []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", filePath, err)
		}
	}

	// Test GetTopDirectories with count=3
	topDirs, err := GetTopDirectories(tempDir, 3)
	if err != nil {
		t.Fatalf("GetTopDirectories failed: %v", err)
	}

	// Check the results
	if len(topDirs) != 3 {
		t.Errorf("len(topDirs) = %d, expected 3", len(topDirs))
	}

	// Check that the directories are sorted by size in descending order
	for i := 0; i < len(topDirs)-1; i++ {
		if topDirs[i].Size < topDirs[i+1].Size {
			t.Errorf("topDirs[%d].Size = %d, topDirs[%d].Size = %d, expected descending order",
				i, topDirs[i].Size, i+1, topDirs[i+1].Size)
		}
	}

	// Test GetTopDirectories with count=0 (should return all directories)
	allDirs, err := GetTopDirectories(tempDir, 0)
	if err != nil {
		t.Fatalf("GetTopDirectories with count=0 failed: %v", err)
	}

	// Check the results
	expectedDirCount := len(dirs)
	if len(allDirs) != expectedDirCount {
		t.Errorf("len(allDirs) = %d, expected %d", len(allDirs), expectedDirCount)
	}
}
