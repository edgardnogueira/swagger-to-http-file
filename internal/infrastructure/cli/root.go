package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "swagger-to-http-file",
	Short: "Convert Swagger/OpenAPI documents to .http files",
	Long: `A CLI tool that converts Swagger/OpenAPI JSON documents into .http files 
	for easy API testing. It handles various parameter types and 
	can organize requests by tags.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Future flag definitions will go here
}
