package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// SizeInfo represents information about the size of a Docker image
type SizeInfo struct {
	TotalSize         int64            `json:"totalSize"`
	LayerSizes        []LayerSize      `json:"layerSizes,omitempty"`
	LargestFiles      []FileSize       `json:"largestFiles,omitempty"`
	LargestDirs       []DirectoryInfo  `json:"largestDirs,omitempty"`
	FileTypeBreakdown map[string]int64 `json:"fileTypeBreakdown,omitempty"`
}

// LayerSize represents the size of a Docker image layer
type LayerSize struct {
	ID   string `json:"id"`
	Size int64  `json:"size"`
}

// FileSize represents the size of a file
type FileSize struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// AnalyzeSize analyzes the size of a Docker image
func AnalyzeSize(imagePath string) (*SizeInfo, error) {
	info := &SizeInfo{
		FileTypeBreakdown: make(map[string]int64),
	}

	// Get the total size of the image
	totalSize, err := getDirSize(imagePath)
	if err != nil {
		return nil, err
	}
	info.TotalSize = totalSize

	// Get the size of each layer
	layerSizes, err := getLayerSizes(imagePath)
	if err != nil {
		return nil, err
	}
	info.LayerSizes = layerSizes

	// Get the largest files
	largestFiles, err := getLargestFiles(imagePath, 10)
	if err != nil {
		return nil, err
	}
	info.LargestFiles = largestFiles

	// Get the largest directories
	largestDirs, err := GetTopDirectories(imagePath, 5)
	if err != nil {
		return nil, err
	}
	info.LargestDirs = largestDirs

	// Get the file type breakdown
	fileTypeBreakdown, err := getFileTypeBreakdown(imagePath)
	if err != nil {
		return nil, err
	}
	info.FileTypeBreakdown = fileTypeBreakdown

	return info, nil
}

// getLayerSizes returns the size of each layer in a Docker image
func getLayerSizes(imagePath string) ([]LayerSize, error) {
	var layerSizes []LayerSize

	// Find all layer directories
	layerDirs, err := filepath.Glob(filepath.Join(imagePath, "*/layer.tar"))
	if err != nil {
		return nil, fmt.Errorf("failed to find layer tarballs: %v", err)
	}

	// Get the size of each layer
	for _, layerPath := range layerDirs {
		// Get the layer ID from the directory name
		layerDir := filepath.Dir(layerPath)
		layerID := filepath.Base(layerDir)

		// Get the size of the layer tarball
		info, err := os.Stat(layerPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get layer size: %v", err)
		}

		// Add the layer size to the list
		layerSizes = append(layerSizes, LayerSize{
			ID:   layerID,
			Size: info.Size(),
		})
	}

	// Sort the layers by size in descending order
	sort.Slice(layerSizes, func(i, j int) bool {
		return layerSizes[i].Size > layerSizes[j].Size
	})

	return layerSizes, nil
}

// getLargestFiles returns the N largest files in a Docker image
func getLargestFiles(imagePath string, count int) ([]FileSize, error) {
	var files []FileSize

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

		// Add the file to the list
		files = append(files, FileSize{
			Path: relPath,
			Size: info.Size(),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort the files by size in descending order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	// Return the top N files
	if count > 0 && count < len(files) {
		return files[:count], nil
	}
	return files, nil
}

// getFileTypeBreakdown returns the breakdown of file types by size
func getFileTypeBreakdown(imagePath string) (map[string]int64, error) {
	fileTypeBreakdown := make(map[string]int64)

	// Walk the directory tree
	err := filepath.Walk(imagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get the file extension
		ext := filepath.Ext(path)
		if ext == "" {
			ext = "[no extension]"
		}

		// Add the file size to the breakdown
		fileTypeBreakdown[ext] += info.Size()

		return nil
	})

	return fileTypeBreakdown, err
}

// FormatSize formats a size in bytes to a human-readable string
func FormatSize(size int64) string {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
	)

	var (
		unit  string
		value float64
	)

	switch {
	case size >= int64(TB):
		unit = "TB"
		value = float64(size) / TB
	case size >= int64(GB):
		unit = "GB"
		value = float64(size) / GB
	case size >= int64(MB):
		unit = "MB"
		value = float64(size) / MB
	case size >= int64(KB):
		unit = "KB"
		value = float64(size) / KB
	default:
		unit = "B"
		value = float64(size)
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}
