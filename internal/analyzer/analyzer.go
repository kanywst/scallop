package analyzer

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// AnalysisResult represents the result of a Docker image analysis
type AnalysisResult struct {
	ImagePath     string          `json:"imagePath"`
	AnalyzedAt    time.Time       `json:"analyzedAt"`
	DirectoryInfo *DirectoryInfo  `json:"directoryInfo,omitempty"`
	SecurityInfo  *SecurityResult `json:"securityInfo,omitempty"`
	SizeInfo      *SizeInfo       `json:"sizeInfo,omitempty"`
}

// AnalyzeImage analyzes a Docker image
func AnalyzeImage(imagePath string, verbose bool) *AnalysisResult {
	result := &AnalysisResult{
		ImagePath:  imagePath,
		AnalyzedAt: time.Now(),
	}

	// Analyze directory structure
	fmt.Println("Analyzing directory structure...")
	dirInfo, err := AnalyzeDirectory(imagePath, verbose)
	if err != nil {
		fmt.Printf("Error analyzing directory structure: %v\n", err)
	} else {
		result.DirectoryInfo = dirInfo
	}

	// Analyze security
	fmt.Println("Analyzing security...")
	securityInfo, err := AnalyzeSecurity(imagePath)
	if err != nil {
		fmt.Printf("Error analyzing security: %v\n", err)
	} else {
		result.SecurityInfo = securityInfo
	}

	// Analyze size
	fmt.Println("Analyzing size...")
	sizeInfo, err := AnalyzeSize(imagePath)
	if err != nil {
		fmt.Printf("Error analyzing size: %v\n", err)
	} else {
		result.SizeInfo = sizeInfo
	}

	return result
}

// OutputJSON outputs the analysis result as JSON
func OutputJSON(result *AnalysisResult, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// OutputText outputs the analysis result as text
func OutputText(result *AnalysisResult, writer io.Writer) error {
	// Output header
	fmt.Fprintf(writer, "Docker Image Analysis: %s\n", result.ImagePath)
	fmt.Fprintf(writer, "Analyzed at: %s\n\n", result.AnalyzedAt.Format(time.RFC1123))

	// Output directory info
	if result.DirectoryInfo != nil {
		fmt.Fprintf(writer, "Directory Structure:\n")
		fmt.Fprintf(writer, "  Files: %d\n", result.DirectoryInfo.FileCount)
		fmt.Fprintf(writer, "  Directories: %d\n", result.DirectoryInfo.DirCount)
		fmt.Fprintf(writer, "  Total Size: %s\n", FormatSize(result.DirectoryInfo.Size))

		// Output file types
		if len(result.DirectoryInfo.FileTypes) > 0 {
			fmt.Fprintf(writer, "  File Types:\n")
			for ext, count := range result.DirectoryInfo.FileTypes {
				fmt.Fprintf(writer, "    %s: %d files\n", ext, count)
			}
		}
		fmt.Fprintln(writer)
	}

	// Output security info
	if result.SecurityInfo != nil {
		fmt.Fprintf(writer, "Security Analysis:\n")
		fmt.Fprintf(writer, "  Total Issues: %d\n", result.SecurityInfo.TotalIssues)
		fmt.Fprintf(writer, "  High Severity: %d\n", result.SecurityInfo.HighSeverity)
		fmt.Fprintf(writer, "  Medium Severity: %d\n", result.SecurityInfo.MediumSeverity)
		fmt.Fprintf(writer, "  Low Severity: %d\n", result.SecurityInfo.LowSeverity)

		// Output issues
		if len(result.SecurityInfo.Issues) > 0 {
			fmt.Fprintf(writer, "  Issues:\n")
			for _, issue := range result.SecurityInfo.Issues {
				fmt.Fprintf(writer, "    [%s] %s: %s (%s)\n", issue.Severity, issue.Type, issue.Description, issue.Path)
			}
		}
		fmt.Fprintln(writer)
	}

	// Output size info
	if result.SizeInfo != nil {
		fmt.Fprintf(writer, "Size Analysis:\n")
		fmt.Fprintf(writer, "  Total Size: %s\n", FormatSize(result.SizeInfo.TotalSize))

		// Output layer sizes
		if len(result.SizeInfo.LayerSizes) > 0 {
			fmt.Fprintf(writer, "  Layer Sizes:\n")
			for _, layer := range result.SizeInfo.LayerSizes {
				fmt.Fprintf(writer, "    %s: %s\n", layer.ID, FormatSize(layer.Size))
			}
		}

		// Output largest files
		if len(result.SizeInfo.LargestFiles) > 0 {
			fmt.Fprintf(writer, "  Largest Files:\n")
			for _, file := range result.SizeInfo.LargestFiles {
				fmt.Fprintf(writer, "    %s: %s\n", file.Path, FormatSize(file.Size))
			}
		}

		// Output largest directories
		if len(result.SizeInfo.LargestDirs) > 0 {
			fmt.Fprintf(writer, "  Largest Directories:\n")
			for _, dir := range result.SizeInfo.LargestDirs {
				fmt.Fprintf(writer, "    %s: %s\n", dir.Path, FormatSize(dir.Size))
			}
		}

		// Output file type breakdown
		if len(result.SizeInfo.FileTypeBreakdown) > 0 {
			fmt.Fprintf(writer, "  File Type Breakdown:\n")
			for ext, size := range result.SizeInfo.FileTypeBreakdown {
				fmt.Fprintf(writer, "    %s: %s\n", ext, FormatSize(size))
			}
		}
		fmt.Fprintln(writer)
	}

	// Output recommendations
	fmt.Fprintf(writer, "Recommendations:\n")
	recommendations := generateRecommendations(result)
	for _, rec := range recommendations {
		fmt.Fprintf(writer, "  - %s\n", rec)
	}

	return nil
}

// generateRecommendations generates recommendations based on the analysis result
func generateRecommendations(result *AnalysisResult) []string {
	var recommendations []string

	// Security recommendations
	if result.SecurityInfo != nil {
		if result.SecurityInfo.HighSeverity > 0 {
			recommendations = append(recommendations, fmt.Sprintf("Address %d high severity security issues", result.SecurityInfo.HighSeverity))
		}
		if result.SecurityInfo.MediumSeverity > 0 {
			recommendations = append(recommendations, fmt.Sprintf("Review %d medium severity security issues", result.SecurityInfo.MediumSeverity))
		}

		// Check for specific issue types
		sensitiveFiles := 0
		hardcodedSecrets := 0
		vulnPackages := 0
		for _, issue := range result.SecurityInfo.Issues {
			switch issue.Type {
			case "SENSITIVE_FILE":
				sensitiveFiles++
			case "HARDCODED_SECRET":
				hardcodedSecrets++
			case "VULNERABLE_PACKAGE":
				vulnPackages++
			}
		}

		if sensitiveFiles > 0 {
			recommendations = append(recommendations, fmt.Sprintf("Remove %d sensitive files from the image", sensitiveFiles))
		}
		if hardcodedSecrets > 0 {
			recommendations = append(recommendations, fmt.Sprintf("Remove %d hardcoded secrets from the image", hardcodedSecrets))
		}
		if vulnPackages > 0 {
			recommendations = append(recommendations, fmt.Sprintf("Update %d vulnerable packages", vulnPackages))
		}
	}

	// Size recommendations
	if result.SizeInfo != nil {
		// Check if the image is large
		if result.SizeInfo.TotalSize > 500*1024*1024 { // 500 MB
			recommendations = append(recommendations, "Consider reducing the image size")

			// Check if there are large files
			if len(result.SizeInfo.LargestFiles) > 0 && result.SizeInfo.LargestFiles[0].Size > 100*1024*1024 { // 100 MB
				recommendations = append(recommendations, fmt.Sprintf("Remove or optimize large files (largest: %s, %s)", result.SizeInfo.LargestFiles[0].Path, FormatSize(result.SizeInfo.LargestFiles[0].Size)))
			}

			// Check for specific file types that might be unnecessary
			unnecessaryExts := []string{".zip", ".tar", ".gz", ".log", ".tmp"}
			for _, ext := range unnecessaryExts {
				if size, ok := result.SizeInfo.FileTypeBreakdown[ext]; ok && size > 10*1024*1024 { // 10 MB
					recommendations = append(recommendations, fmt.Sprintf("Remove unnecessary %s files (%s)", ext, FormatSize(size)))
				}
			}
		}

		// Check if there are many layers
		if len(result.SizeInfo.LayerSizes) > 10 {
			recommendations = append(recommendations, "Reduce the number of layers in the image")
		}
	}

	// Add general recommendations if none were generated
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "No specific recommendations at this time")
	}

	return recommendations
}
