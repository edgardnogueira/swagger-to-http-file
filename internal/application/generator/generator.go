package generator

import (
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// HttpGenerator defines the interface for generating HTTP files from Swagger documents
type HttpGenerator interface {
	Generate(doc *models.SwaggerDoc) (map[string]*models.HttpFile, error)
}

// Basic implementation will be added during issue implementation
