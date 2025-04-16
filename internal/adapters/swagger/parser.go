package swagger

import (
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// Parser implements the application.SwaggerParser interface
type Parser struct {}

// Parse parses Swagger JSON data and returns a SwaggerDoc
func (p *Parser) Parse(data []byte) (*models.SwaggerDoc, error) {
	// Will be implemented as part of Issue #1
	return nil, nil
}

// Validate validates the SwaggerDoc structure
func (p *Parser) Validate(doc *models.SwaggerDoc) error {
	// Will be implemented as part of Issue #1
	return nil
}
