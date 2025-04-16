package parser

import (
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// SwaggerParser defines the interface for parsing Swagger documents
type SwaggerParser interface {
	Parse(data []byte) (*models.SwaggerDoc, error)
	Validate(doc *models.SwaggerDoc) error
}

// Basic implementation will be added during issue implementation
