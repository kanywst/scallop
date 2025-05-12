package analyzer

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DirectoryInfo represents information about a directory
type DirectoryInfo struct {
	Path      string         `json:"path"`
	Size      int64          `json:"size"`
	FileCount int            `json:"fileCount"`
	DirCount  int            `json:"dirCount"`
	Files     []string       `json:"files,omitempty"`
	FileTypes map[string]int `json:"fileTypes,omitempty"`
	Dirs      []string       `json:"dirs,omitempty"`
}

// AnalyzeDirectory analyzes the directory structure of a Docker image
func AnalyzeDirectory(imagePath string, verbose bool) (*DirectoryInfo, error) {
	info := &DirectoryInfo{
		Path:      imagePath,
		FileTypes: make(map[string]int),
	}

	// Walk the directory tree
	err := filepath.Walk(imagePath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == imagePath {
			return nil
		}

		// Get the relative path
		relPath, err := filepath.Rel(imagePath, path)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			// Count directories
			info.DirCount++
			if verbose {
				info.Dirs = append(info.Dirs, relPath)
			}
		} else {
			// Count files
			info.FileCount++
			info.Size += fileInfo.Size()

			// Count file types
			ext := strings.ToLower(filepath.Ext(path))
			if ext == "" {
				ext = "[no extension]"
			}
			info.FileTypes[ext]++

			if verbose {
				info.Files = append(info.Files, relPath)
			}
		}

		return nil
	})

	// Sort the files and directories for better readability
	if verbose {
		sort.Strings(info.Files)
		sort.Strings(info.Dirs)
	}

	return info, err
}

// GetTopDirectories returns the top N directories by size
func GetTopDirectories(imagePath string, count int) ([]DirectoryInfo, error) {
	var dirs []DirectoryInfo

	// Get all directories in the image
	err := filepath.Walk(imagePath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory and non-directories
		if path == imagePath || !fileInfo.IsDir() {
			return nil
		}

		// Get the directory size
		size, err := getDirSize(path)
		if err != nil {
			return err
		}

		// Get the relative path
		relPath, err := filepath.Rel(imagePath, path)
		if err != nil {
			return err
		}

		// Add the directory to the list
		dirs = append(dirs, DirectoryInfo{
			Path: relPath,
			Size: size,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort the directories by size in descending order
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Size > dirs[j].Size
	})

	// Return the top N directories
	if count > 0 && count < len(dirs) {
		return dirs[:count], nil
	}
	return dirs, nil
}

// getDirSize returns the total size of all files in a directory in bytes
func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
