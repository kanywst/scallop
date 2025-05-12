package main

import (
	"fmt"
	"os"

	"github.com/kanywst/scallop/internal/analyzer"
	"github.com/kanywst/scallop/internal/docker"
)

func main() {
	// Check if an image name is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run example.go <image_name>")
		fmt.Println("Example: go run example.go nginx:latest")
		os.Exit(1)
	}

	imageName := os.Args[1]
	fmt.Printf("Analyzing Docker image: %s\n", imageName)

	// Create a temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "scallop-example-")
	if err != nil {
		fmt.Printf("Error creating temporary directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	// Extract the Docker image
	fmt.Println("Extracting Docker image...")
	extractedPath, err := docker.ExtractImage(imageName, tempDir)
	if err != nil {
		fmt.Printf("Error extracting Docker image: %v\n", err)
		os.Exit(1)
	}

	// Analyze the extracted image
	fmt.Println("Analyzing Docker image...")
	results := analyzer.AnalyzeImage(extractedPath, true)

	// Output the results as text
	fmt.Println("\nAnalysis Results:")
	analyzer.OutputText(results, os.Stdout)
}
