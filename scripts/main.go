package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	upstreamURL     = "https://raw.githubusercontent.com/HodlBook/hodlbook/master/umbrel/umbrel-app.yml"
	placeholderText = "PLACEHOLDER_VERSION"
)

func main() {
	projectRoot, _ := os.Getwd()
	placeholderDir := filepath.Join(projectRoot, "placeholder")
	hodlbookDir := filepath.Join(projectRoot, "hodlbook")

	if _, err := os.Stat(placeholderDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: placeholder directory not found at %s\n", placeholderDir)
		os.Exit(1)
	}

	fmt.Println("Fetching latest version from HodlBook/hodlbook...")

	version, err := fetchVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching version: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Latest version: %s\n\n", version)

	if err := os.MkdirAll(hodlbookDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	files := []string{"umbrel-app.yml", "docker-compose.yml"}
	for _, file := range files {
		if err := processTemplate(placeholderDir, hodlbookDir, file, version); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", file, err)
			os.Exit(1)
		}
	}

	fmt.Println("\nSync completed successfully!")
}

func fetchVersion() (string, error) {
	resp, err := http.Get(upstreamURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Extract version using regex: version: "X.Y.Z" or version: X.Y.Z
	re := regexp.MustCompile(`version:\s*"?([0-9]+\.[0-9]+\.[0-9]+)"?`)
	matches := re.FindSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("version not found in upstream file")
	}

	return string(matches[1]), nil
}

func processTemplate(srcDir, destDir, filename, version string) error {
	srcPath := filepath.Join(srcDir, filename)
	destPath := filepath.Join(destDir, filename)

	fmt.Printf("Processing %s...\n", filename)

	content, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Replace placeholder with actual version
	// For docker-compose.yml, use vX.Y.Z format for image tag
	newContent := string(content)
	if filename == "docker-compose.yml" {
		newContent = strings.ReplaceAll(newContent, placeholderText, "v"+version)
	} else {
		newContent = strings.ReplaceAll(newContent, placeholderText, version)
	}

	if err := os.WriteFile(destPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("  ✓ Written to %s (version: %s)\n", destPath, version)
	return nil
}
