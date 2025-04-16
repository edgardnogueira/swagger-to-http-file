package http

import (
	"fmt"
	"strings"

	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// Formatter handles formatting HTTP files into the .http format
type Formatter struct{}

// FormatHttpFile formats an HttpFile into a string representation in .http format
func (f *Formatter) FormatHttpFile(file *models.HttpFile) string {
	var builder strings.Builder

	// Add base URL and global variables
	builder.WriteString(f.formatGlobalVars(file.GlobalVars))
	builder.WriteString("\n")

	// Add requests
	for i, req := range file.Requests {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(f.FormatHttpRequest(req))
		builder.WriteString("\n")
	}

	return builder.String()
}

// FormatHttpRequest formats a single HttpRequest into a string representation
func (f *Formatter) FormatHttpRequest(req models.HttpRequest) string {
	var builder strings.Builder

	// Add request name as a comment
	builder.WriteString(fmt.Sprintf("### %s\n", req.Name))

	// Add description if present
	if req.Description != "" {
		builder.WriteString(fmt.Sprintf("# %s\n", req.Description))
	}

	// Add method and URL
	builder.WriteString(fmt.Sprintf("%s {{baseUrl}}%s\n", req.Method, req.Path))

	// Add headers
	for name, value := range req.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\n", name, value))
	}

	// Add body if present
	if req.Body != "" {
		builder.WriteString("\n")
		builder.WriteString(req.Body)
		builder.WriteString("\n")
	}

	return builder.String()
}

// formatGlobalVars formats global variables for the .http file
func (f *Formatter) formatGlobalVars(vars map[string]string) string {
	var builder strings.Builder

	// Add comment header for variables
	if len(vars) > 0 {
		builder.WriteString("# Global variables\n")
	}

	// Add each variable
	for name, value := range vars {
		builder.WriteString(fmt.Sprintf("@%s = %s\n", name, value))
	}

	return builder.String()
}

// formatNamedRequest formats a request with a name and response capture
func (f *Formatter) formatNamedRequest(req models.HttpRequest) string {
	var builder strings.Builder

	// Add request name as a comment
	builder.WriteString(fmt.Sprintf("### %s\n", req.Name))

	// Add request name for response capture
	builder.WriteString(fmt.Sprintf("# @name %s\n", getVarName(req.Name)))

	// Add description if present
	if req.Description != "" {
		builder.WriteString(fmt.Sprintf("# %s\n", req.Description))
	}

	// Add method and URL
	builder.WriteString(fmt.Sprintf("%s {{baseUrl}}%s\n", req.Method, req.Path))

	// Add headers
	for name, value := range req.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\n", name, value))
	}

	// Add body if present
	if req.Body != "" {
		builder.WriteString("\n")
		builder.WriteString(req.Body)
		builder.WriteString("\n")
	}

	return builder.String()
}

// getVarName converts a request name to a valid variable name
func getVarName(name string) string {
	// Replace spaces with underscores and remove special characters
	varName := strings.ToLower(name)
	varName = strings.ReplaceAll(varName, " ", "_")
	varName = strings.ReplaceAll(varName, "-", "_")

	// Remove any non-alphanumeric characters
	var result strings.Builder
	for i, char := range varName {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_' {
			// First character must be a letter or underscore
			if i == 0 && char >= '0' && char <= '9' {
				result.WriteRune('_')
			}
			result.WriteRune(char)
		}
	}

	return result.String()
}

// NewFormatter creates a new Formatter instance
func NewFormatter() *Formatter {
	return &Formatter{}
}
