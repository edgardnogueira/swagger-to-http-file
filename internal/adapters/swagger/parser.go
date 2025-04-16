package swagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// Parser implements the application.SwaggerParser interface
type Parser struct{}

// Parse parses Swagger JSON data and returns a SwaggerDoc
func (p *Parser) Parse(data []byte) (*models.SwaggerDoc, error) {
	var doc models.SwaggerDoc
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse Swagger JSON: %w", err)
	}

	return &doc, nil
}

// Validate validates the SwaggerDoc structure
func (p *Parser) Validate(doc *models.SwaggerDoc) error {
	if doc == nil {
		return errors.New("swagger document is nil")
	}

	// Check if it's a valid Swagger/OpenAPI document
	if doc.Swagger == "" && doc.OpenAPI == "" {
		return errors.New("not a valid Swagger/OpenAPI document: missing 'swagger' or 'openapi' field")
	}

	// Validate required fields
	if doc.Info.Title == "" {
		return errors.New("info.title is required")
	}

	if doc.Info.Version == "" {
		return errors.New("info.version is required")
	}

	// Check if there are any paths
	if len(doc.Paths) == 0 {
		return errors.New("no paths defined in the document")
	}

	// Validate paths and operations
	if err := p.validatePaths(doc); err != nil {
		return err
	}

	return nil
}

// validatePaths validates the paths and operations in the Swagger document
func (p *Parser) validatePaths(doc *models.SwaggerDoc) error {
	for path, pathItem := range doc.Paths {
		if !strings.HasPrefix(path, "/") {
			return fmt.Errorf("path %q must begin with a forward slash", path)
		}
		
		if err := p.validatePathItem(path, &pathItem); err != nil {
			return err
		}
	}
	return nil
}

// validatePathItem validates a single path item and its operations
func (p *Parser) validatePathItem(path string, item *models.PathItem) error {
	// Validate operations
	if item.Get != nil {
		if err := p.validateOperation(path, "GET", item.Get); err != nil {
			return err
		}
	}
	if item.Post != nil {
		if err := p.validateOperation(path, "POST", item.Post); err != nil {
			return err
		}
	}
	if item.Put != nil {
		if err := p.validateOperation(path, "PUT", item.Put); err != nil {
			return err
		}
	}
	if item.Delete != nil {
		if err := p.validateOperation(path, "DELETE", item.Delete); err != nil {
			return err
		}
	}
	if item.Options != nil {
		if err := p.validateOperation(path, "OPTIONS", item.Options); err != nil {
			return err
		}
	}
	if item.Head != nil {
		if err := p.validateOperation(path, "HEAD", item.Head); err != nil {
			return err
		}
	}
	if item.Patch != nil {
		if err := p.validateOperation(path, "PATCH", item.Patch); err != nil {
			return err
		}
	}
	return nil
}

// validateOperation validates a single operation in a path
func (p *Parser) validateOperation(path, method string, op *models.Operation) error {
	// Check if responses are defined (required by OpenAPI spec)
	if len(op.Responses) == 0 {
		return fmt.Errorf("no responses defined for %s %s", method, path)
	}

	// Validate parameters
	for i, param := range op.Parameters {
		if param.Name == "" {
			return fmt.Errorf("parameter %d for %s %s has no name", i, method, path)
		}
		
		if param.In == "" {
			return fmt.Errorf("parameter '%s' for %s %s has no 'in' property", param.Name, method, path)
		}
		
		// Path parameters must be required
		if param.In == "path" && !param.Required {
			return fmt.Errorf("path parameter '%s' for %s %s must be required", param.Name, method, path)
		}
	}

	return nil
}

// GetBaseURL computes the base URL from the Swagger document
func (p *Parser) GetBaseURL(doc *models.SwaggerDoc) string {
	// For OpenAPI v3
	if len(doc.Servers) > 0 && doc.Servers[0].URL != "" {
		return doc.Servers[0].URL
	}
	
	// For Swagger v2
	if doc.Schemes != nil && len(doc.Schemes) > 0 {
		return fmt.Sprintf("%s://%s%s", doc.Schemes[0], doc.Host, doc.BasePath)
	}
	
	if doc.Host != "" {
		return fmt.Sprintf("http://%s%s", doc.Host, doc.BasePath)
	}
	
	// Default fallback
	return "http://localhost"
}

// ExtractOperations extracts all operations from the Swagger document
func (p *Parser) ExtractOperations(doc *models.SwaggerDoc) map[string][]models.OperationInfo {
	operations := make(map[string][]models.OperationInfo)
	
	for path, pathItem := range doc.Paths {
		p.addOperation(operations, path, "GET", pathItem.Get)
		p.addOperation(operations, path, "POST", pathItem.Post)
		p.addOperation(operations, path, "PUT", pathItem.Put)
		p.addOperation(operations, path, "DELETE", pathItem.Delete)
		p.addOperation(operations, path, "OPTIONS", pathItem.Options)
		p.addOperation(operations, path, "HEAD", pathItem.Head)
		p.addOperation(operations, path, "PATCH", pathItem.Patch)
	}
	
	return operations
}

// addOperation adds an operation to the operations map, organized by tag
func (p *Parser) addOperation(operations map[string][]models.OperationInfo, path, method string, op *models.Operation) {
	if op == nil {
		return
	}
	
	info := models.OperationInfo{
		Path:        path,
		Method:      method,
		Operation:   op,
		Parameters:  op.Parameters,
	}
	
	// Group by tag, or use "default" if no tags present
	tags := op.Tags
	if len(tags) == 0 {
		tags = []string{"default"}
	}
	
	for _, tag := range tags {
		operations[tag] = append(operations[tag], info)
	}
}

// New creates a new Parser instance
func New() *Parser {
	return &Parser{}
}
