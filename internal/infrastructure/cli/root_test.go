package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHelperFunctions(t *testing.T) {
	// Test fileExists
	t.Run("fileExists", func(t *testing.T) {
		// Create temp file
		tempFile, err := os.CreateTemp("", "test-file")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		// Test existing file
		if !fileExists(tempFile.Name()) {
			t.Errorf("Expected file %s to exist", tempFile.Name())
		}

		// Test non-existing file
		if fileExists(tempFile.Name() + ".nonexistent") {
			t.Errorf("Expected file %s to not exist", tempFile.Name()+".nonexistent")
		}
	})

	// Test dirExists
	t.Run("dirExists", func(t *testing.T) {
		// Test existing directory
		if !dirExists(os.TempDir()) {
			t.Errorf("Expected directory %s to exist", os.TempDir())
		}

		// Test non-existing directory
		if dirExists(filepath.Join(os.TempDir(), "nonexistent-dir-"+randString(8))) {
			t.Errorf("Expected directory to not exist")
		}
	})

	// Test sanitizeFilename
	t.Run("sanitizeFilename", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "simple name",
				input:    "simple",
				expected: "simple",
			},
			{
				name:     "with path traversal",
				input:    "../../../etc/passwd",
				expected: "etc/passwd",
			},
			{
				name:     "clean dots",
				input:    "file....txt",
				expected: "file....txt",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := sanitizeFilename(test.input)
				if result != test.expected {
					t.Errorf("Expected sanitized filename to be %s, got %s", test.expected, result)
				}
			})
		}
	})
}

// Helper functions for tests
func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
