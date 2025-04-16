package http

import (
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// Generator implements the application.HttpGenerator interface
type Generator struct{}

// Generate creates HTTP files from a Swagger document
func (g *Generator) Generate(doc *models.SwaggerDoc) (map[string]*models.HttpFile, error) {
	// Will be implemented as part of Issue #2
	return nil, nil
}
