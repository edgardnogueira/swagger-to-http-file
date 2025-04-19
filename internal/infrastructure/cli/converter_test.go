package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/edgardnogueira/swagger-to-http-file/internal/adapters/http"
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

func TestSanitizeTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple tag",
			input:    "pets",
			expected: "pets",
		},
		{
			name:     "tag with spaces",
			input:    "Pet Store",
			expected: "pet_store",
		},
		{
			name:     "tag with special characters",
			input:    "Users/Auth",
			expected: "users_auth",
		},
		{
			name:     "tag with mixed case",
			input:    "UserAuthentication",
			expected: "userauthentication",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := sanitizeTag(test.input)
			if result != test.expected {
				t.Errorf("Expected sanitized tag to be %s, got %s", test.expected, result)
			}
		})
	}
}

func TestExtractGlobalVars(t *testing.T) {
	// Create test files
	files := map[string]*models.HTTPFile{
		"tag1": {
			GlobalVars: map[string]string{
				"baseUrl": "http://api1.example.com",
				"apiKey":  "key1",
			},
		},
		"tag2": {
			GlobalVars: map[string]string{
				"baseUrl":  "http://api2.example.com",
				"username": "user1",
			},
		},
	}

	// Extract global vars
	vars := extractGlobalVars(files)

	// Check results
	expectedKeys := []string{"baseUrl", "apiKey", "username"}
	for _, key := range expectedKeys {
		if _, exists := vars[key]; !exists {
			t.Errorf("Expected global var '%s' not found", key)
		}
	}

	// Check specific values (last one wins for duplicates)
	if vars["baseUrl"] != "http://api2.example.com" {
		t.Errorf("Expected baseUrl to be http://api2.example.com, got %s", vars["baseUrl"])
	}
}

func TestWriteHTTPFiles(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "swagger-to-http-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	HTTPFiles := map[string]*models.HTTPFile{
		"pets": {
			BaseURL: "http://api.example.com",
			GlobalVars: map[string]string{
				"baseUrl": "http://api.example.com",
				"apiKey":  "test-key",
			},
			Requests: []models.HTTPRequest{
				{
					Name:   "Get Pets",
					Method: "GET",
					Path:   "/pets",
				},
			},
			Tag: "pets",
		},
		"users": {
			BaseURL: "http://api.example.com",
			GlobalVars: map[string]string{
				"baseUrl": "http://api.example.com",
			},
			Requests: []models.HTTPRequest{
				{
					Name:   "Get Users",
					Method: "GET",
					Path:   "/users",
				},
			},
			Tag: "users",
		},
	}

	// Create formatter
	formatter := http.NewFormatter()

	// Test group by tag
	t.Run("group by tag", func(t *testing.T) {
		dir := filepath.Join(tempDir, "grouped")
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		err = WriteHTTPFiles(HTTPFiles, dir, formatter, true, true, false)
		if err != nil {
			t.Fatalf("WriteHTTPFiles failed: %v", err)
		}

		// Check if files were created
		expectedFiles := []string{"pets.http", "users.http"}
		for _, file := range expectedFiles {
			path := filepath.Join(dir, file)
			if !fileExists(path) {
				t.Errorf("Expected file %s to be created", path)
			}

			// Check content
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", path, err)
				continue
			}

			// Simple content checks
			if len(content) == 0 {
				t.Errorf("File %s is empty", path)
			}
		}
	})

	// Test single file
	t.Run("single file", func(t *testing.T) {
		dir := filepath.Join(tempDir, "single")
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		err = WriteHTTPFiles(HTTPFiles, dir, formatter, false, true, false)
		if err != nil {
			t.Fatalf("WriteHTTPFiles failed: %v", err)
		}

		// Check if file was created
		path := filepath.Join(dir, "swagger.http")
		if !fileExists(path) {
			t.Errorf("Expected file %s to be created", path)
		}

		// Check content
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", path, err)
		} else if len(content) == 0 {
			t.Errorf("File %s is empty", path)
		}
	})

	// Test overwrite flag
	t.Run("overwrite flag", func(t *testing.T) {
		dir := filepath.Join(tempDir, "overwrite")
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Create existing file
		path := filepath.Join(dir, "pets.http")
		err = os.WriteFile(path, []byte("existing content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// First run with overwrite=false
		err = WriteHTTPFiles(HTTPFiles, dir, formatter, true, false, false)
		if err != nil {
			t.Fatalf("WriteHTTPFiles failed: %v", err)
		}

		// Check content (should still be original)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", path, err)
		} else if string(content) != "existing content" {
			t.Errorf("Expected file content to remain unchanged")
		}

		// Run with overwrite=true
		err = WriteHTTPFiles(HTTPFiles, dir, formatter, true, true, false)
		if err != nil {
			t.Fatalf("WriteHTTPFiles failed: %v", err)
		}

		// Check content (should be updated)
		content, err = os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", path, err)
		} else if string(content) == "existing content" {
			t.Errorf("Expected file content to be updated")
		}
	})
}
