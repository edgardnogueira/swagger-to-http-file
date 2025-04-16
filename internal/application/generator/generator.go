package generator

import (
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// HttpGenerator defines the interface for generating HTTP files from Swagger documents
type HttpGenerator interface {
	// Generate creates HTTP files from a Swagger document
	Generate(doc *models.SwaggerDoc, baseURL string) (map[string]*models.HttpFile, error)

	// GenerateRequest creates a single HTTP request from an operation
	GenerateRequest(op models.OperationInfo, baseURL string) models.HttpRequest

	// FormatPath formats path parameters for use in a request URL
	FormatPath(path string, params []models.Parameter) string

	// ExtractGlobalVars extracts global variables that could be used across requests
	ExtractGlobalVars(doc *models.SwaggerDoc) map[string]string
}
