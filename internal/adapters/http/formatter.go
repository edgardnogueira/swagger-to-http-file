package http

import (
	"fmt"
	"regexp"
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

/*
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
*/
func getVarName(name string) string {
	// 1) normalize spaces & hyphens to underscores
	s := strings.ReplaceAll(name, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")

	// 2) insert underscore between a lowercase/digit and uppercase (camelCase boundary)
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	s = re.ReplaceAllString(s, `${1}_${2}`)

	// 3) lowercase everything
	s = strings.ToLower(s)

	// 4) strip out any remaining non-alphanumeric/_, and ensure it doesn't start with a digit
	var result strings.Builder
	for i, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' {
			if i == 0 && (ch >= '0' && ch <= '9') {
				result.WriteRune('_')
			}
			result.WriteRune(ch)
		}
	}

	return result.String()
}

// NewFormatter creates a new Formatter instance
func NewFormatter() *Formatter {
	return &Formatter{}
}
