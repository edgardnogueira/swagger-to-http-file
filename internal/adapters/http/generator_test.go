package http

import (
	"os"
	"testing"

	"github.com/edgardnogueira/swagger-to-http-file/internal/adapters/swagger"
	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

func TestGenerator_Generate(t *testing.T) {
	// Load test data
	data, err := os.ReadFile("../../../test/samples/petstore.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Set up parser and generator
	parser := swagger.New()
	generator := New(parser)

	// Parse Swagger document
	doc, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("Failed to parse swagger: %v", err)
	}

	// Generate HTTP files
	baseURL := parser.GetBaseURL(doc)
	files, err := generator.Generate(doc, baseURL)
	if err != nil {
		t.Fatalf("Failed to generate HTTP files: %v", err)
	}

	// Verify files were generated correctly
	if len(files) == 0 {
		t.Errorf("No HTTP files were generated")
	}

	// Check if "pets" tag exists
	petsFile, exists := files["pets"]
	if !exists {
		t.Errorf("Expected 'pets' file to be generated")
	} else {
		// Verify the pets file
		if len(petsFile.Requests) == 0 {
			t.Errorf("No requests in the pets file")
		}

		// Check if baseURL is set
		if petsFile.BaseURL == "" {
			t.Errorf("Base URL not set in pets file")
		}

		// Check for global variables
		if len(petsFile.GlobalVars) == 0 {
			t.Errorf("No global variables in pets file")
		}

		// Inspect a few requests
		hasGetRequest := false
		hasPostRequest := false
		for _, req := range petsFile.Requests {
			if req.Method == "GET" && req.Path == "/pets" {
				hasGetRequest = true
			}
			if req.Method == "POST" && req.Path == "/pets" {
				hasPostRequest = true
			}
		}

		if !hasGetRequest {
			t.Errorf("Missing GET /pets request")
		}
		if !hasPostRequest {
			t.Errorf("Missing POST /pets request")
		}
	}
}

func TestGenerator_GenerateRequest(t *testing.T) {
	parser := swagger.New()
	generator := New(parser)
	baseURL := "http://example.com/api"

	// Test with a GET operation
	getOp := models.OperationInfo{
		Path:   "/pets",
		Method: "GET",
		Operation: &models.Operation{
			Summary:     "List all pets",
			OperationID: "listPets",
			Tags:        []string{"pets"},
			Parameters: []models.Parameter{
				{
					Name:        "limit",
					In:          "query",
					Description: "How many items to return",
					Required:    false,
					Type:        "integer",
				},
			},
			Responses: map[string]models.Response{
				"200": {
					Description: "A list of pets",
				},
			},
		},
		Parameters: []models.Parameter{
			{
				Name:        "limit",
				In:          "query",
				Description: "How many items to return",
				Required:    false,
				Type:        "integer",
			},
		},
	}

	// Test with a POST operation with body
	postOp := models.OperationInfo{
		Path:   "/pets",
		Method: "POST",
		Operation: &models.Operation{
			Summary:     "Create a pet",
			OperationID: "createPet",
			Tags:        []string{"pets"},
			Parameters: []models.Parameter{
				{
					Name:        "pet",
					In:          "body",
					Description: "Pet to create",
					Required:    true,
					Schema: &models.SchemaObj{
						Type: "object",
						Properties: map[string]models.SchemaObj{
							"name": {
								Type: "string",
							},
							"age": {
								Type: "integer",
							},
						},
					},
				},
			},
			Responses: map[string]models.Response{
				"201": {
					Description: "Pet created",
				},
			},
		},
		Parameters: []models.Parameter{
			{
				Name:        "pet",
				In:          "body",
				Description: "Pet to create",
				Required:    true,
				Schema: &models.SchemaObj{
					Type: "object",
					Properties: map[string]models.SchemaObj{
						"name": {
							Type: "string",
						},
						"age": {
							Type: "integer",
						},
					},
				},
			},
		},
	}

	// Generate and check GET request
	getReq := generator.GenerateRequest(getOp, baseURL)
	if getReq.Method != "GET" {
		t.Errorf("Expected GET method, got %s", getReq.Method)
	}
	if getReq.Path != "/pets" {
		t.Errorf("Expected /pets path, got %s", getReq.Path)
	}
	if getReq.Name != "List all pets" {
		t.Errorf("Expected 'List all pets' name, got %s", getReq.Name)
	}
	if len(getReq.Body) > 0 {
		t.Errorf("GET request should not have a body")
	}

	// Generate and check POST request
	postReq := generator.GenerateRequest(postOp, baseURL)
	if postReq.Method != "POST" {
		t.Errorf("Expected POST method, got %s", postReq.Method)
	}
	if postReq.Path != "/pets" {
		t.Errorf("Expected /pets path, got %s", postReq.Path)
	}
	if postReq.Name != "Create a pet" {
		t.Errorf("Expected 'Create a pet' name, got %s", postReq.Name)
	}
	if len(postReq.Body) == 0 {
		t.Errorf("POST request should have a body")
	}
	if contentType, exists := postReq.Headers["Content-Type"]; !exists || contentType != "application/json" {
		t.Errorf("Expected Content-Type header to be application/json")
	}
}

func TestGenerator_FormatPath(t *testing.T) {
	parser := swagger.New()
	generator := New(parser)

	tests := []struct {
		name     string
		path     string
		params   []models.Parameter
		expected string
	}{
		{
			name:     "no parameters",
			path:     "/pets",
			params:   []models.Parameter{},
			expected: "/pets",
		},
		{
			name: "with path parameter",
			path: "/pets/{petId}",
			params: []models.Parameter{
				{
					Name: "petId",
					In:   "path",
				},
			},
			expected: "/pets/{{petId}}",
		},
		{
			name: "with multiple path parameters",
			path: "/users/{userId}/pets/{petId}",
			params: []models.Parameter{
				{
					Name: "userId",
					In:   "path",
				},
				{
					Name: "petId",
					In:   "path",
				},
			},
			expected: "/users/{{userId}}/pets/{{petId}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.FormatPath(tt.path, tt.params)
			if result != tt.expected {
				t.Errorf("FormatPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerator_ExtractGlobalVars(t *testing.T) {
	parser := swagger.New()
	generator := New(parser)

	tests := []struct {
		name           string
		doc            *models.SwaggerDoc
		expectedVars   []string
		expectedValues map[string]string
	}{
		{
			name: "with host and basePath",
			doc: &models.SwaggerDoc{
				Host:     "api.example.com",
				BasePath: "/v1",
				Schemes:  []string{"https"},
			},
			expectedVars: []string{"baseUrl", "authToken"},
			expectedValues: map[string]string{
				"baseUrl": "https://api.example.com/v1",
			},
		},
		{
			name: "with servers (OpenAPI v3)",
			doc: &models.SwaggerDoc{
				Servers: []models.Server{
					{
						URL: "https://api.example.com/v2",
					},
				},
			},
			expectedVars: []string{"baseUrl", "authToken"},
			expectedValues: map[string]string{
				"baseUrl": "https://api.example.com/v2",
			},
		},
		{
			name:         "with minimal info",
			doc:          &models.SwaggerDoc{},
			expectedVars: []string{"baseUrl", "authToken"},
			expectedValues: map[string]string{
				"baseUrl": "http://localhost",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := generator.ExtractGlobalVars(tt.doc)

			// Check that all expected vars exist
			for _, expectedVar := range tt.expectedVars {
				if _, exists := vars[expectedVar]; !exists {
					t.Errorf("Expected variable %s not found", expectedVar)
				}
			}

			// Check specific values
			for k, expectedValue := range tt.expectedValues {
				if value, exists := vars[k]; !exists || value != expectedValue {
					t.Errorf("Expected %s=%s, got %s", k, expectedValue, value)
				}
			}
		})
	}
}
