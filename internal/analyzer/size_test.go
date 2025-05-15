package analyzer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnalyzeSize(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "size-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a directory structure for testing
	dirs := []string{
		"dir1",
		"dir2",
		"dir3",
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
		{"dir3/file7.jpg", "test content 7"},
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

	// Create mock layer directories and files
	layerDirs := []string{
		"layer1",
		"layer2",
		"layer3",
	}

	for _, dir := range layerDirs {
		dirPath := filepath.Join(tempDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create layer directory %q: %v", dirPath, err)
		}

		// Create a layer.tar file in each layer directory
		layerPath := filepath.Join(dirPath, "layer.tar")
		if err := os.WriteFile(layerPath, []byte("mock layer content "+dir), 0644); err != nil {
			t.Fatalf("Failed to create layer file %q: %v", layerPath, err)
		}
	}

	// Run the size analysis
	result, err := AnalyzeSize(tempDir)
	if err != nil {
		t.Fatalf("AnalyzeSize failed: %v", err)
	}

	// Check the results
	if result.TotalSize <= 0 {
		t.Errorf("TotalSize = %d, expected > 0", result.TotalSize)
	}

	// Check layer sizes
	if len(result.LayerSizes) != len(layerDirs) {
		t.Errorf("len(LayerSizes) = %d, expected %d", len(result.LayerSizes), len(layerDirs))
	}

	// Check that layers are sorted by size in descending order
	for i := 0; i < len(result.LayerSizes)-1; i++ {
		if result.LayerSizes[i].Size < result.LayerSizes[i+1].Size {
			t.Errorf("LayerSizes[%d].Size = %d, LayerSizes[%d].Size = %d, expected descending order",
				i, result.LayerSizes[i].Size, i+1, result.LayerSizes[i+1].Size)
		}
	}

	// Check largest files
	if len(result.LargestFiles) == 0 {
		t.Errorf("LargestFiles is empty, expected some files")
	}

	// Check largest directories
	if len(result.LargestDirs) == 0 {
		t.Errorf("LargestDirs is empty, expected some directories")
	}

	// Check file type breakdown
	if len(result.FileTypeBreakdown) == 0 {
		t.Errorf("FileTypeBreakdown is empty, expected some file types")
	}

	// Check specific file types
	expectedFileTypes := []string{".txt", ".go", ".jpg"}
	for _, ext := range expectedFileTypes {
		if result.FileTypeBreakdown[ext] <= 0 {
			t.Errorf("FileTypeBreakdown[%q] = %d, expected > 0", ext, result.FileTypeBreakdown[ext])
		}
	}
}

func TestGetLayerSizes(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "layer-sizes-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create mock layer directories and files with different sizes
	layerDirs := []struct {
		name    string
		content string
	}{
		{"layer1", "small content"},
		{"layer2", "medium content with more bytes"},
		{"layer3", "large content with even more bytes than the previous one"},
	}

	for _, layer := range layerDirs {
		dirPath := filepath.Join(tempDir, layer.name)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create layer directory %q: %v", dirPath, err)
		}

		// Create a layer.tar file in each layer directory
		layerPath := filepath.Join(dirPath, "layer.tar")
		if err := os.WriteFile(layerPath, []byte(layer.content), 0644); err != nil {
			t.Fatalf("Failed to create layer file %q: %v", layerPath, err)
		}
	}

	// Run the layer sizes check
	layerSizes, err := getLayerSizes(tempDir)
	if err != nil {
		t.Fatalf("getLayerSizes failed: %v", err)
	}

	// Check the results
	if len(layerSizes) != len(layerDirs) {
		t.Errorf("len(layerSizes) = %d, expected %d", len(layerSizes), len(layerDirs))
	}

	// Check that layers are sorted by size in descending order
	for i := 0; i < len(layerSizes)-1; i++ {
		if layerSizes[i].Size < layerSizes[i+1].Size {
			t.Errorf("layerSizes[%d].Size = %d, layerSizes[%d].Size = %d, expected descending order",
				i, layerSizes[i].Size, i+1, layerSizes[i+1].Size)
		}
	}

	// Check that the largest layer is layer3
	if layerSizes[0].ID != "layer3" {
		t.Errorf("layerSizes[0].ID = %q, expected %q", layerSizes[0].ID, "layer3")
	}

	// Check that the smallest layer is layer1
	if layerSizes[len(layerSizes)-1].ID != "layer1" {
		t.Errorf("layerSizes[%d].ID = %q, expected %q", len(layerSizes)-1, layerSizes[len(layerSizes)-1].ID, "layer1")
	}
}

func TestGetLargestFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "largest-files-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create files with different sizes
	files := []struct {
		path    string
		content string
	}{
		{"small.txt", "small content"},
		{"medium.txt", "medium content with more bytes"},
		{"large.txt", "large content with even more bytes than the previous one"},
		{"dir/nested.txt", "nested file with some content"},
		{"dir/subdir/deep.txt", "deeply nested file with some content"},
	}

	// Create files
	for _, file := range files {
		dirPath := filepath.Dir(filepath.Join(tempDir, file.path))
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %q: %v", dirPath, err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, file.path), []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", file.path, err)
		}
	}

	// Run the largest files check with count=3
	largestFiles, err := getLargestFiles(tempDir, 3)
	if err != nil {
		t.Fatalf("getLargestFiles failed: %v", err)
	}

	// Check the results
	if len(largestFiles) != 3 {
		t.Errorf("len(largestFiles) = %d, expected 3", len(largestFiles))
	}

	// Check that files are sorted by size in descending order
	for i := 0; i < len(largestFiles)-1; i++ {
		if largestFiles[i].Size < largestFiles[i+1].Size {
			t.Errorf("largestFiles[%d].Size = %d, largestFiles[%d].Size = %d, expected descending order",
				i, largestFiles[i].Size, i+1, largestFiles[i+1].Size)
		}
	}

	// Check that the largest file is large.txt
	if !strings.HasSuffix(largestFiles[0].Path, "large.txt") {
		t.Errorf("largestFiles[0].Path = %q, expected to end with %q", largestFiles[0].Path, "large.txt")
	}

	// Run the largest files check with count=0 (should return all files)
	allFiles, err := getLargestFiles(tempDir, 0)
	if err != nil {
		t.Fatalf("getLargestFiles with count=0 failed: %v", err)
	}

	// Check the results
	if len(allFiles) != len(files) {
		t.Errorf("len(allFiles) = %d, expected %d", len(allFiles), len(files))
	}
}

func TestGetFileTypeBreakdown(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file-type-breakdown-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create files with different extensions
	files := []struct {
		path    string
		content string
	}{
		{"file1.txt", "text file content"},
		{"file2.txt", "another text file"},
		{"file3.go", "go file content"},
		{"file4.go", "another go file"},
		{"file5.jpg", "mock jpg content"},
		{"file6", "file without extension"},
	}

	// Create files
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(tempDir, file.path), []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", file.path, err)
		}
	}

	// Run the file type breakdown check
	breakdown, err := getFileTypeBreakdown(tempDir)
	if err != nil {
		t.Fatalf("getFileTypeBreakdown failed: %v", err)
	}

	// Check the results
	expectedTypes := map[string]int{
		".txt":           2,
		".go":            2,
		".jpg":           1,
		"[no extension]": 1,
	}

	for ext, _ := range expectedTypes {
		size := breakdown[ext]
		if size <= 0 {
			t.Errorf("breakdown[%q] = %d, expected > 0", ext, size)
		}

		// Calculate expected size for this extension
		var expectedSize int64
		for _, file := range files {
			fileExt := filepath.Ext(file.path)
			if fileExt == "" && ext == "[no extension]" {
				expectedSize += int64(len(file.content))
			} else if fileExt == ext {
				expectedSize += int64(len(file.content))
			}
		}

		if breakdown[ext] != expectedSize {
			t.Errorf("breakdown[%q] = %d, expected %d", ext, breakdown[ext], expectedSize)
		}
	}
}

// TestFormatSizeInSize is already tested in analyzer_test.go
