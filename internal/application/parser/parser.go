package parser

import (
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// SwaggerParser defines the interface for parsing Swagger documents
type SwaggerParser interface {
	// Parse parses raw Swagger data into a structured document
	Parse(data []byte) (*models.SwaggerDoc, error)
	
	// Validate checks if the Swagger document is valid
	Validate(doc *models.SwaggerDoc) error
	
	// GetBaseURL computes the base URL from the Swagger document
	GetBaseURL(doc *models.SwaggerDoc) string
	
	// ExtractOperations extracts all operations from the document, grouped by tag
	ExtractOperations(doc *models.SwaggerDoc) map[string][]models.OperationInfo
}
