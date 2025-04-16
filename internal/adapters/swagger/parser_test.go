package swagger

import (
	"os"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	// Load the test Swagger file
	data, err := os.ReadFile("../../../test/samples/petstore.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Create a new parser
	parser := New()

	// Parse the Swagger document
	doc, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("Failed to parse Swagger document: %v", err)
	}

	// Verify basic document structure
	if doc.Swagger != "2.0" {
		t.Errorf("Expected swagger version 2.0, got %s", doc.Swagger)
	}

	if doc.Info.Title != "Swagger Petstore" {
		t.Errorf("Expected title 'Swagger Petstore', got %s", doc.Info.Title)
	}

	if doc.Host != "petstore.swagger.io" {
		t.Errorf("Expected host 'petstore.swagger.io', got %s", doc.Host)
	}

	// Verify paths
	if len(doc.Paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(doc.Paths))
	}

	// Check if /pets path exists
	petsPath, exists := doc.Paths["/pets"]
	if !exists {
		t.Errorf("Expected /pets path to exist")
	} else {
		// Verify GET operation
		if petsPath.Get == nil {
			t.Errorf("Expected GET operation on /pets path")
		} else {
			if petsPath.Get.Summary != "List all pets" {
				t.Errorf("Expected GET summary 'List all pets', got %s", petsPath.Get.Summary)
			}
		}

		// Verify POST operation
		if petsPath.Post == nil {
			t.Errorf("Expected POST operation on /pets path")
		} else {
			if petsPath.Post.Summary != "Create a pet" {
				t.Errorf("Expected POST summary 'Create a pet', got %s", petsPath.Post.Summary)
			}
		}
	}

	// Check if /pets/{petId} path exists
	petPath, exists := doc.Paths["/pets/{petId}"]
	if !exists {
		t.Errorf("Expected /pets/{petId} path to exist")
	} else {
		// Verify GET operation
		if petPath.Get == nil {
			t.Errorf("Expected GET operation on /pets/{petId} path")
		} else {
			if petPath.Get.Summary != "Info for a specific pet" {
				t.Errorf("Expected GET summary 'Info for a specific pet', got %s", petPath.Get.Summary)
			}
		}
	}
}

func TestParser_Validate(t *testing.T) {
	// Load the test Swagger file
	data, err := os.ReadFile("../../../test/samples/petstore.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Create a new parser
	parser := New()

	// Parse the Swagger document
	doc, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("Failed to parse Swagger document: %v", err)
	}

	// Validate the document
	err = parser.Validate(doc)
	if err != nil {
		t.Errorf("Valid Swagger document failed validation: %v", err)
	}

	// Test validation with nil document
	err = parser.Validate(nil)
	if err == nil {
		t.Errorf("Expected error for nil document")
	}

	// Test validation with invalid document
	invalidDoc := *doc
	invalidDoc.Swagger = ""
	invalidDoc.OpenAPI = ""
	err = parser.Validate(&invalidDoc)
	if err == nil {
		t.Errorf("Expected error for invalid document without swagger/openapi version")
	}

	// Test validation with missing info title
	invalidDoc = *doc
	invalidDoc.Info.Title = ""
	err = parser.Validate(&invalidDoc)
	if err == nil {
		t.Errorf("Expected error for document with missing info.title")
	}
}

func TestParser_GetBaseURL(t *testing.T) {
	// Load the test Swagger file
	data, err := os.ReadFile("../../../test/samples/petstore.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Create a new parser
	parser := New()

	// Parse the Swagger document
	doc, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("Failed to parse Swagger document: %v", err)
	}

	// Get the base URL
	baseURL := parser.GetBaseURL(doc)
	expected := "http://petstore.swagger.io/api"
	if baseURL != expected {
		t.Errorf("Expected base URL %s, got %s", expected, baseURL)
	}

	// Test with OpenAPI v3 servers
	v3Doc := *doc
	v3Doc.Swagger = ""
	v3Doc.OpenAPI = "3.0.0"
	v3Doc.Host = ""
	v3Doc.BasePath = ""
	v3Doc.Schemes = nil
	v3Doc.Servers = []struct {
		URL         string                    `json:"url"`
		Description string                    `json:"description,omitempty"`
		Variables   map[string]struct {
			Enum        []string `json:"enum,omitempty"`
			Default     string   `json:"default"`
			Description string   `json:"description,omitempty"`
		} `json:"variables,omitempty"`
	}{
		{
			URL:         "https://api.example.com/v1",
			Description: "Production server",
		},
	}

	baseURL = parser.GetBaseURL(&v3Doc)
	expected = "https://api.example.com/v1"
	if baseURL != expected {
		t.Errorf("Expected base URL %s, got %s", expected, baseURL)
	}

	// Test with fallback
	v3Doc.Servers = nil
	baseURL = parser.GetBaseURL(&v3Doc)
	expected = "http://localhost"
	if baseURL != expected {
		t.Errorf("Expected fallback base URL %s, got %s", expected, baseURL)
	}
}

func TestParser_ExtractOperations(t *testing.T) {
	// Load the test Swagger file
	data, err := os.ReadFile("../../../test/samples/petstore.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Create a new parser
	parser := New()

	// Parse the Swagger document
	doc, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("Failed to parse Swagger document: %v", err)
	}

	// Extract operations
	operations := parser.ExtractOperations(doc)

	// Check that the pets tag exists
	petsOps, exists := operations["pets"]
	if !exists {
		t.Errorf("Expected 'pets' tag to exist in operations")
	}

	// Check that we have the correct number of operations
	if len(petsOps) != 5 {
		t.Errorf("Expected 5 operations under 'pets' tag, got %d", len(petsOps))
	}

	// Verify the operations
	methodsFound := map[string]bool{
		"GET /pets":          false,
		"POST /pets":         false,
		"GET /pets/{petId}":  false,
		"PUT /pets/{petId}":  false,
		"DELETE /pets/{petId}": false,
	}

	for _, op := range petsOps {
		key := op.Method + " " + op.Path
		if _, ok := methodsFound[key]; ok {
			methodsFound[key] = true
		} else {
			t.Errorf("Unexpected operation: %s", key)
		}
	}

	for op, found := range methodsFound {
		if !found {
			t.Errorf("Expected operation %s not found", op)
		}
	}
}
