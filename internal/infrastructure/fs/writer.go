package fs

import (
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// FileWriter handles writing HTTP files to the filesystem
type FileWriter struct {
	OutputDir string
}

// WriteHTTPFiles writes the HTTP files to the filesystem
func (fw *FileWriter) WriteHTTPFiles(files map[string]*models.HTTPFile) error {
	// Will be implemented as part of Issue #3
	return nil
}
