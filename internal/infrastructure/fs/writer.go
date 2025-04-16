package fs

import (
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

// FileWriter handles writing HTTP files to the filesystem
type FileWriter struct {
	OutputDir string
}

// WriteHttpFiles writes the HTTP files to the filesystem
func (fw *FileWriter) WriteHttpFiles(files map[string]*models.HttpFile) error {
	// Will be implemented as part of Issue #3
	return nil
}
