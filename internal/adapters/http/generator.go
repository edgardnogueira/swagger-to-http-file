package http

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/edgardnogueira/swagger-to-http-file/internal/application/parser"
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// Generator implements the application.HttpGenerator interface
type Generator struct {
	parser parser.SwaggerParser
}

// Generate creates HTTP files from a Swagger document
func (g *Generator) Generate(doc *models.SwaggerDoc, baseURL string) (map[string]*models.HttpFile, error) {
	if doc == nil {
		return nil, fmt.Errorf("swagger document is nil")
	}

	// Extract global variables
	globalVars := g.ExtractGlobalVars(doc)

	// Extract operations by tag
	operations := g.parser.ExtractOperations(doc)

	// Create HTTP files per tag
	files := make(map[string]*models.HttpFile)

	for tag, ops := range operations {
		httpFile := &models.HttpFile{
			BaseURL:    baseURL,
			GlobalVars: globalVars,
			Requests:   []models.HttpRequest{},
			Tag:        tag,
		}

		// Generate requests for each operation
		for _, op := range ops {
			request := g.GenerateRequest(op, baseURL)
			httpFile.Requests = append(httpFile.Requests, request)
		}

		files[tag] = httpFile
	}

	return files, nil
}

// GenerateRequest creates a single HTTP request from an operation
func (g *Generator) GenerateRequest(op models.OperationInfo, baseURL string) models.HttpRequest {
	// Format path with parameters
	path := g.FormatPath(op.Path, op.Parameters)

	// Create request
	request := models.HttpRequest{
		Name:        generateRequestName(op),
		Method:      op.Method,
		Path:        path,
		Headers:     extractHeaders(op),
		Body:        generateRequestBody(op),
		Description: generateDescription(op),
		Vars:        extractVars(op),
		Tag:         getFirstTag(op.Operation),
	}

	return request
}

// FormatPath formats path parameters for use in a request URL
func (g *Generator) FormatPath(path string, params []models.Parameter) string {
	// For .http files, path parameters are referenced as {{paramName}}
	re := regexp.MustCompile(`\{([^}]+)\}`)
	return re.ReplaceAllString(path, "{{$1}}")
}

// ExtractGlobalVars extracts global variables that could be used across requests
func (g *Generator) ExtractGlobalVars(doc *models.SwaggerDoc) map[string]string {
	vars := make(map[string]string)

	// Add baseUrl variable
	if len(doc.Servers) > 0 && doc.Servers[0].URL != "" {
		vars["baseUrl"] = doc.Servers[0].URL
	} else if doc.Host != "" {
		scheme := "http"
		if len(doc.Schemes) > 0 {
			scheme = doc.Schemes[0]
		}
		vars["baseUrl"] = fmt.Sprintf("%s://%s%s", scheme, doc.Host, doc.BasePath)
	} else {
		vars["baseUrl"] = "http://localhost"
	}

	// Add common auth variables
	vars["authToken"] = "your_auth_token"

	return vars
}

// Helper functions for request generation

// generateRequestName creates a readable name for the request
func generateRequestName(op models.OperationInfo) string {
	if op.Operation.Summary != "" {
		return op.Operation.Summary
	}

	if op.Operation.OperationID != "" {
		return toTitleCase(op.Operation.OperationID)
	}

	// Fallback to method and path
	return fmt.Sprintf("%s %s", op.Method, op.Path)
}

// toTitleCase converts camelCase or snake_case to Title Case
func toTitleCase(s string) string {
	// Replace underscores with spaces
	s = strings.ReplaceAll(s, "_", " ")

	// Insert spaces before capital letters
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	s = re.ReplaceAllString(s, "$1 $2")

	// Title case the entire string
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[0:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// extractHeaders extracts headers from the operation
func extractHeaders(op models.OperationInfo) map[string]string {
	headers := make(map[string]string)

	// Add Content-Type header based on consumes
	if op.Operation.Consumes != nil && len(op.Operation.Consumes) > 0 {
		headers["Content-Type"] = op.Operation.Consumes[0]
	} else {
		// Default to JSON for operations with request bodies
		for _, param := range op.Parameters {
			if param.In == "body" || isBodyParameter(param) {
				headers["Content-Type"] = "application/json"
				break
			}
		}
	}

	// Add Accept header based on produces
	if op.Operation.Produces != nil && len(op.Operation.Produces) > 0 {
		headers["Accept"] = op.Operation.Produces[0]
	}

	// Add Authorization header if security is defined
	if op.Operation.Security != nil && len(op.Operation.Security) > 0 {
		// Generic placeholder for auth token
		headers["Authorization"] = "Bearer {{authToken}}"
	}

	// Extract header parameters
	for _, param := range op.Parameters {
		if param.In == "header" {
			// Use parameter name as-is (headers are case-insensitive)
			headers[param.Name] = fmt.Sprintf("{{%s}}", param.Name)
		}
	}

	return headers
}

// isBodyParameter checks if a parameter is a body parameter in OpenAPI v3
func isBodyParameter(param models.Parameter) bool {
	return param.In == "body" ||
		(param.Schema != nil && param.Schema.Type == "object") ||
		(param.In == "query" && param.Style == "form" && param.Explode)
}

// generateRequestBody generates a request body example based on the operation
func generateRequestBody(op models.OperationInfo) string {
	// Look for body parameters
	for _, param := range op.Parameters {
		if param.In == "body" && param.Schema != nil {
			return generateSchemaExample(param.Schema)
		}
	}

	// Check for request body (OpenAPI v3)
	if op.Operation.RequestBody != nil && op.Operation.RequestBody.Content != nil {
		for contentType, mediaType := range op.Operation.RequestBody.Content {
			if strings.Contains(contentType, "json") && mediaType.Schema != nil {
				return generateSchemaExample(mediaType.Schema)
			}
		}
	}

	return ""
}

// generateSchemaExample generates an example JSON for a schema
func generateSchemaExample(schema *models.SchemaObj) string {
	// For $ref schemas, we can't resolve them without a full document
	if schema.Ref != "" {
		return "{\n  // Reference to " + schema.Ref + "\n  // Replace with actual data\n}"
	}

	// Handle primitive types
	switch schema.Type {
	case "string":
		if schema.Example != nil {
			return fmt.Sprintf("\"%v\"", schema.Example)
		}
		return "\"string\""
	case "integer", "number":
		if schema.Example != nil {
			return fmt.Sprintf("%v", schema.Example)
		}
		return "0"
	case "boolean":
		if schema.Example != nil {
			return fmt.Sprintf("%v", schema.Example)
		}
		return "false"
	case "array":
		if schema.Items != nil {
			itemExample := generateSchemaExample(schema.Items)
			return fmt.Sprintf("[\n  %s\n]", itemExample)
		}
		return "[]"
	case "object":
		var builder strings.Builder
		builder.WriteString("{\n")

		if schema.Properties != nil {
			i := 0
			for propName, propSchema := range schema.Properties {
				propExample := generateSchemaExample(&propSchema)
				if i > 0 {
					builder.WriteString(",\n")
				}
				builder.WriteString(fmt.Sprintf("  \"%s\": %s", propName, propExample))
				i++
			}
		}

		builder.WriteString("\n}")
		return builder.String()
	default:
		return "{}"
	}
}

// generateDescription generates a description for the request
func generateDescription(op models.OperationInfo) string {
	var desc strings.Builder

	if op.Operation.Description != "" {
		desc.WriteString(op.Operation.Description)
	} else if op.Operation.Summary != "" {
		desc.WriteString(op.Operation.Summary)
	}

	return desc.String()
}

// extractVars extracts variables from the operation
func extractVars(op models.OperationInfo) map[string]string {
	vars := make(map[string]string)

	// Add path parameters as variables
	for _, param := range op.Parameters {
		if param.In == "path" {
			// Generate example value based on type
			varValue := generateExampleValue(param)
			vars[param.Name] = varValue
		}
	}

	return vars
}

// generateExampleValue generates an example value for a parameter
func generateExampleValue(param models.Parameter) string {
	// Use example if provided
	if param.Example != nil {
		return fmt.Sprintf("%v", param.Example)
	}

	// Use default if provided
	if param.Default != nil {
		return fmt.Sprintf("%v", param.Default)
	}

	// Generate based on type
	switch param.Type {
	case "string":
		return "example_string"
	case "integer", "number":
		return "123"
	case "boolean":
		return "true"
	default:
		return "example"
	}
}

// getFirstTag gets the first tag of an operation or returns "default"
func getFirstTag(op *models.Operation) string {
	if op.Tags != nil && len(op.Tags) > 0 {
		return op.Tags[0]
	}
	return "default"
}

// New creates a new Generator instance
func New(p parser.SwaggerParser) *Generator {
	return &Generator{
		parser: p,
	}
}
