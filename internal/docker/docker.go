package docker

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kanywst/scallop/internal/utils"
)

// IsDockerImageName checks if the provided string is a Docker image name rather than a file path
func IsDockerImageName(name string) bool {
	return !strings.Contains(name, "/") || strings.Contains(name, ":")
}

// ExtractImage extracts a Docker image to the specified directory
// It handles both local tar files and Docker image names from Docker daemon
func ExtractImage(imagePath string, destDir string) (string, error) {
	if IsDockerImageName(imagePath) {
		return extractFromDockerDaemon(imagePath, destDir)
	}
	return extractFromTarFile(imagePath, destDir)
}

// extractFromDockerDaemon saves a Docker image from the Docker daemon and extracts it
func extractFromDockerDaemon(imageName string, destDir string) (string, error) {
	// Create a temporary file to save the Docker image
	tempFile, err := os.CreateTemp("", "docker-image-*.tar")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Save the Docker image to a tar file
	saveCmd := exec.Command("docker", "save", "-o", tempFile.Name(), imageName)
	if err := saveCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to save Docker image: %v", err)
	}

	// Extract the saved image
	return extractFromTarFile(tempFile.Name(), destDir)
}

// extractFromTarFile extracts a Docker image from a tar file
func extractFromTarFile(tarPath string, destDir string) (string, error) {
	// Open the tar file
	file, err := os.Open(tarPath)
	if err != nil {
		return "", fmt.Errorf("failed to open tar file: %v", err)
	}
	defer file.Close()

	// Create the destination directory for the extracted image
	imageDir := filepath.Join(destDir, "image")
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create image directory: %v", err)
	}

	// Check if it's a gzipped tar file
	var tarReader *tar.Reader
	if strings.HasSuffix(tarPath, ".gz") || strings.HasSuffix(tarPath, ".tgz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return "", fmt.Errorf("failed to create gzip reader: %v", err)
		}
		defer gzReader.Close()
		tarReader = tar.NewReader(gzReader)
	} else {
		tarReader = tar.NewReader(file)
	}

	// Extract the tar file
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading tar file: %v", err)
		}

		// Skip if the header is nil
		if header == nil {
			continue
		}

		// Create the file path
		target := filepath.Join(imageDir, header.Name)

		// Check for path traversal attacks
		if !strings.HasPrefix(target, imageDir) {
			return "", fmt.Errorf("invalid tar file: contains path traversal attack")
		}

		// Handle different types of files
		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(target, 0755); err != nil {
				return "", fmt.Errorf("failed to create directory: %v", err)
			}
		case tar.TypeReg:
			// Create directory for the file if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return "", fmt.Errorf("failed to create directory: %v", err)
			}

			// Create the file
			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return "", fmt.Errorf("failed to create file: %v", err)
			}

			// Copy the file content
			if _, err := io.Copy(file, tarReader); err != nil {
				file.Close()
				return "", fmt.Errorf("failed to copy file content: %v", err)
			}
			file.Close()
		case tar.TypeSymlink:
			// Create symlink
			if err := os.Symlink(header.Linkname, target); err != nil {
				return "", fmt.Errorf("failed to create symlink: %v", err)
			}
		default:
			// Skip other types of files
		}
	}

	// Extract layer tarballs if they exist
	if err := extractLayers(imageDir); err != nil {
		return "", fmt.Errorf("failed to extract layers: %v", err)
	}

	return imageDir, nil
}

// extractLayers extracts the layer tarballs in a Docker image
func extractLayers(imageDir string) error {
	// Find all layer tarballs
	layerDirs, err := filepath.Glob(filepath.Join(imageDir, "*/layer.tar"))
	if err != nil {
		return fmt.Errorf("failed to find layer tarballs: %v", err)
	}

	// Extract each layer tarball
	for _, layerPath := range layerDirs {
		layerDir := filepath.Dir(layerPath)
		extractedDir := filepath.Join(layerDir, "extracted")

		// Create the extracted directory
		if err := os.MkdirAll(extractedDir, 0755); err != nil {
			return fmt.Errorf("failed to create extracted directory: %v", err)
		}

		// Extract the layer tarball
		if err := utils.ExtractTar(layerPath, extractedDir); err != nil {
			return fmt.Errorf("failed to extract layer tarball: %v", err)
		}
	}

	return nil
}
