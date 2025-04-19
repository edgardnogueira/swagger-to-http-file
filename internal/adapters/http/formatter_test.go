package http

import (
	"strings"
	"testing"

	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

func TestFormatter_FormatHTTPRequest(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		name     string
		request  models.HTTPRequest
		expected []string
	}{
		{
			name: "basic GET request",
			request: models.HTTPRequest{
				Name:   "Get Pets",
				Method: "GET",
				Path:   "/pets",
				Headers: map[string]string{
					"Accept": "application/json",
				},
				Description: "Get a list of pets",
			},
			expected: []string{
				"### Get Pets",
				"# Get a list of pets",
				"GET {{baseUrl}}/pets",
				"Accept: application/json",
			},
		},
		{
			name: "POST request with body",
			request: models.HTTPRequest{
				Name:   "Create Pet",
				Method: "POST",
				Path:   "/pets",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body:        "{\n  \"name\": \"Fluffy\",\n  \"age\": 3\n}",
				Description: "Create a new pet",
			},
			expected: []string{
				"### Create Pet",
				"# Create a new pet",
				"POST {{baseUrl}}/pets",
				"Content-Type: application/json",
				"{\n  \"name\": \"Fluffy\",\n  \"age\": 3\n}",
			},
		},
		{
			name: "request with path parameters",
			request: models.HTTPRequest{
				Name:   "Get Pet by ID",
				Method: "GET",
				Path:   "/pets/{{petId}}",
				Headers: map[string]string{
					"Accept": "application/json",
				},
			},
			expected: []string{
				"### Get Pet by ID",
				"GET {{baseUrl}}/pets/{{petId}}",
				"Accept: application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatHTTPRequest(tt.request)

			// Check that all expected strings are in the result
			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain '%s', but it was not found.\nResult: %s", expected, result)
				}
			}
		})
	}
}

func TestFormatter_FormatHTTPFile(t *testing.T) {
	formatter := NewFormatter()

	file := &models.HTTPFile{
		BaseURL: "http://api.example.com",
		GlobalVars: map[string]string{
			"baseUrl":    "http://api.example.com",
			"authToken":  "your_auth_token",
			"apiVersion": "v1",
		},
		Requests: []models.HTTPRequest{
			{
				Name:   "Get All Items",
				Method: "GET",
				Path:   "/items",
				Headers: map[string]string{
					"Accept": "application/json",
				},
				Description: "Get all items",
			},
			{
				Name:   "Create Item",
				Method: "POST",
				Path:   "/items",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body:        "{\n  \"name\": \"New Item\",\n  \"price\": 19.99\n}",
				Description: "Create a new item",
			},
		},
		Tag: "items",
	}

	result := formatter.FormatHTTPFile(file)

	expectedContent := []string{
		"# Global variables",
		"@baseUrl = http://api.example.com",
		"@authToken = your_auth_token",
		"@apiVersion = v1",
		"### Get All Items",
		"# Get all items",
		"GET {{baseUrl}}/items",
		"Accept: application/json",
		"### Create Item",
		"# Create a new item",
		"POST {{baseUrl}}/items",
		"Content-Type: application/json",
		"{\n  \"name\": \"New Item\",\n  \"price\": 19.99\n}",
	}

	// Check that all expected strings are in the result
	for _, expected := range expectedContent {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected result to contain '%s', but it was not found.\nResult: %s", expected, result)
		}
	}
}

func TestGetVarName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name",
			input:    "Get Pets",
			expected: "get_pets",
		},
		{
			name:     "with special characters",
			input:    "Create Item #123",
			expected: "create_item_123",
		},
		{
			name:     "starts with number",
			input:    "123 Test",
			expected: "_123_test",
		},
		{
			name:     "camelCase",
			input:    "getPetById",
			expected: "get_pet_by_id",
		},
		{
			name:     "with hyphens",
			input:    "test-name-with-hyphens",
			expected: "test_name_with_hyphens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getVarName(tt.input)
			if result != tt.expected {
				t.Errorf("getVarName(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
