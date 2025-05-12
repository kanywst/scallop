package analyzer

import (
	"bytes"
	"testing"
	"time"
)

func TestOutputJSON(t *testing.T) {
	// Create a sample analysis result
	result := &AnalysisResult{
		ImagePath:  "test-image",
		AnalyzedAt: time.Now(),
		DirectoryInfo: &DirectoryInfo{
			Path:      "test-path",
			Size:      1024,
			FileCount: 10,
			DirCount:  5,
			FileTypes: map[string]int{
				".txt": 5,
				".go":  5,
			},
		},
		SecurityInfo: &SecurityResult{
			Issues:         []SecurityIssue{},
			TotalIssues:    0,
			HighSeverity:   0,
			MediumSeverity: 0,
			LowSeverity:    0,
		},
		SizeInfo: &SizeInfo{
			TotalSize:         1024,
			LayerSizes:        []LayerSize{},
			LargestFiles:      []FileSize{},
			LargestDirs:       []DirectoryInfo{},
			FileTypeBreakdown: map[string]int64{},
		},
	}

	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Output the result as JSON
	err := OutputJSON(result, &buf)
	if err != nil {
		t.Fatalf("OutputJSON failed: %v", err)
	}

	// Check if the output is not empty
	if buf.Len() == 0 {
		t.Errorf("OutputJSON produced empty output")
	}

	// Check if the output contains the image path
	if !bytes.Contains(buf.Bytes(), []byte("test-image")) {
		t.Errorf("OutputJSON did not include the image path")
	}
}

func TestOutputText(t *testing.T) {
	// Create a sample analysis result
	result := &AnalysisResult{
		ImagePath:  "test-image",
		AnalyzedAt: time.Now(),
		DirectoryInfo: &DirectoryInfo{
			Path:      "test-path",
			Size:      1024,
			FileCount: 10,
			DirCount:  5,
			FileTypes: map[string]int{
				".txt": 5,
				".go":  5,
			},
		},
		SecurityInfo: &SecurityResult{
			Issues:         []SecurityIssue{},
			TotalIssues:    0,
			HighSeverity:   0,
			MediumSeverity: 0,
			LowSeverity:    0,
		},
		SizeInfo: &SizeInfo{
			TotalSize:         1024,
			LayerSizes:        []LayerSize{},
			LargestFiles:      []FileSize{},
			LargestDirs:       []DirectoryInfo{},
			FileTypeBreakdown: map[string]int64{},
		},
	}

	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Output the result as text
	err := OutputText(result, &buf)
	if err != nil {
		t.Fatalf("OutputText failed: %v", err)
	}

	// Check if the output is not empty
	if buf.Len() == 0 {
		t.Errorf("OutputText produced empty output")
	}

	// Check if the output contains the image path
	if !bytes.Contains(buf.Bytes(), []byte("test-image")) {
		t.Errorf("OutputText did not include the image path")
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0.00 B"},
		{1023, "1023.00 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
		{1024 * 1024 * 1024 * 1024, "1.00 TB"},
	}

	for _, test := range tests {
		result := FormatSize(test.size)
		if result != test.expected {
			t.Errorf("FormatSize(%d) = %s, expected %s", test.size, result, test.expected)
		}
	}
}
