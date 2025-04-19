package cli

import (
	"fmt"
	"io/ioutil"

	"path/filepath"
	"strings"

	"github.com/edgardnogueira/swagger-to-http-file/internal/adapters/http"
	"github.com/edgardnogueira/swagger-to-http-file/internal/adapters/swagger"
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// convertSwaggerToHttp converts a Swagger file to HTTP files
func convertSwaggerToHttp(inputFile, outputDir, baseURLOverride string, groupByTag, overwrite, verbose bool) error {
	// Read the Swagger file
	if verbose {
		fmt.Printf("Reading Swagger file: %s\n", inputFile)
	}

	swaggerData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	// Parse the Swagger file
	if verbose {
		fmt.Println("Parsing Swagger file...")
	}

	parser := swagger.New()
	doc, err := parser.Parse(swaggerData)
	if err != nil {
		return fmt.Errorf("failed to parse Swagger file: %v", err)
	}

	// Validate the Swagger document
	if err := parser.Validate(doc); err != nil {
		return fmt.Errorf("invalid Swagger document: %v", err)
	}

	// Get base URL (use override if provided)
	baseURL := parser.GetBaseURL(doc)
	if baseURLOverride != "" {
		baseURL = baseURLOverride
	}

	if verbose {
		fmt.Printf("Using base URL: %s\n", baseURL)
	}

	// Generate HTTP files
	if verbose {
		fmt.Println("Generating HTTP files...")
	}

	generator := http.New(parser)
	httpFiles, err := generator.Generate(doc, baseURL)
	if err != nil {
		return fmt.Errorf("failed to generate HTTP files: %v", err)
	}

	// Write files to disk
	if verbose {
		fmt.Printf("Writing HTTP files to: %s\n", outputDir)
	}

	formatter := http.NewFormatter()
	return writeHttpFiles(httpFiles, outputDir, formatter, groupByTag, overwrite, verbose)
}

// writeHttpFiles writes the HTTP files to disk
func writeHttpFiles(files map[string]*models.HttpFile, outputDir string, formatter *http.Formatter, groupByTag, overwrite, verbose bool) error {
	if groupByTag {
		// Write each tag to a separate file
		for tag, file := range files {
			filename := sanitizeTag(tag) + ".http"
			fullPath := filepath.Join(outputDir, filename)

			// Check if file exists and overwrite flag is not set
			if fileExists(fullPath) && !overwrite {
				if verbose {
					fmt.Printf("Skipping existing file: %s\n", fullPath)
				}
				continue
			}

			// Format the file content
			content := formatter.FormatHttpFile(file)

			// Write the file
			if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %v", fullPath, err)
			}

			if verbose {
				fmt.Printf("Created HTTP file: %s with %d requests\n", fullPath, len(file.Requests))
			}
		}
	} else {
		// Write all requests to a single file
		combinedFile := &models.HttpFile{
			BaseURL:    "",
			GlobalVars: extractGlobalVars(files),
			Requests:   []models.HttpRequest{},
		}

		// Collect all requests
		for _, file := range files {
			combinedFile.Requests = append(combinedFile.Requests, file.Requests...)
		}

		// Set the base URL from the first file (if any)
		for _, file := range files {
			combinedFile.BaseURL = file.BaseURL
			break
		}

		filename := "swagger.http"
		fullPath := filepath.Join(outputDir, filename)

		// Check if file exists and overwrite flag is not set
		if fileExists(fullPath) && !overwrite {
			if verbose {
				fmt.Printf("Skipping existing file: %s\n", fullPath)
			}
			return nil
		}

		// Format the file content
		content := formatter.FormatHttpFile(combinedFile)

		// Write the file
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", fullPath, err)
		}

		if verbose {
			fmt.Printf("Created HTTP file: %s with %d requests\n", fullPath, len(combinedFile.Requests))
		}
	}

	return nil
}

// extractGlobalVars extracts global variables from all files
func extractGlobalVars(files map[string]*models.HttpFile) map[string]string {
	vars := make(map[string]string)

	// Collect all global variables
	for _, file := range files {
		for k, v := range file.GlobalVars {
			vars[k] = v
		}
	}

	return vars
}

// sanitizeTag makes a tag suitable for use as a filename
func sanitizeTag(tag string) string {
	// Replace spaces with underscores
	tag = strings.ReplaceAll(tag, " ", "_")
	// Remove special characters
	tag = strings.ReplaceAll(tag, "/", "_")
	tag = strings.ReplaceAll(tag, "\\", "_")
	tag = strings.ReplaceAll(tag, ":", "_")
	tag = strings.ReplaceAll(tag, "*", "_")
	tag = strings.ReplaceAll(tag, "?", "_")
	tag = strings.ReplaceAll(tag, "\"", "_")
	tag = strings.ReplaceAll(tag, "<", "_")
	tag = strings.ReplaceAll(tag, ">", "_")
	tag = strings.ReplaceAll(tag, "|", "_")

	// Convert to lowercase for consistency
	tag = strings.ToLower(tag)

	return tag
}
