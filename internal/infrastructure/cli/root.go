package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputDir  string
	baseURL    string
	verbose    bool
	overwrite  bool
	groupByTag bool
)

var rootCmd = &cobra.Command{
	Use:   "swagger-to-http-file",
	Short: "Convert Swagger/OpenAPI documents to .http files",
	Long: `A CLI tool that converts Swagger/OpenAPI JSON documents into .http files 
	for easy API testing. It handles various parameter types and 
	can organize requests by tags.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no args and no inputFile flag, show help
		if inputFile == "" && len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		// If input is provided as argument
		if inputFile == "" && len(args) > 0 {
			inputFile = args[0]
		}

		// Run the main conversion logic
		if err := run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute executes the root command.
// We don't enforce this with cobra to allow for positional argument usage
// Execute executes the root Cobra command.
// It parses flags and runs the conversion if inputs are provided.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Define flags
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Swagger/OpenAPI JSON file to convert (required)")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", ".", "Directory to save .http files")
	rootCmd.PersistentFlags().StringVarP(&baseURL, "baseUrl", "b", "", "Base URL for API requests (overrides the one in Swagger)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&overwrite, "overwrite", "w", false, "Overwrite existing files")
	rootCmd.PersistentFlags().BoolVarP(&groupByTag, "group-by-tag", "g", true, "Group requests by tags into separate files")

	// Make input file required
	// We don't enforce this with cobra to allow for positional argument usage
}

// run is the main function that processes the Swagger file and generates HTTP files
func run() error {
	// Input validation
	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}

	// Check if input file exists
	if !fileExists(inputFile) {
		return fmt.Errorf("input file not found: %s", inputFile)
	}

	// Check if output directory exists, create if not
	if !dirExists(outputDir) {
		if verbose {
			fmt.Printf("Creating output directory: %s\n", outputDir)
		}
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	// This function will be implemented in another file
	if err := convertSwaggerToHTTP(inputFile, outputDir, baseURL, groupByTag, overwrite, verbose); err != nil {
		return err
	}

	return nil
}

// Helper functions
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// sanitizeFilename ensures a safe filename
// sanitizeFilename ensures a safe filename by stripping any “..” or “.” segments
func sanitizeFilename(name string) string {
	// 1) clean the path lexically
	cleaned := filepath.Clean(name)

	// 2) split into path segments
	parts := strings.Split(cleaned, string(filepath.Separator))

	// 3) drop any “.”, “..” or empty segments
	var safeParts []string
	for _, p := range parts {
		if p == "" || p == "." || p == ".." {
			continue
		}
		safeParts = append(safeParts, p)
	}

	// 4) re-join into a sanitized path
	return filepath.Join(safeParts...)
}
