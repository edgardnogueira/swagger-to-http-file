package swagger

import (
	"os"
	"testing"

	"github.com/edgardnogueira/swagger-to-http-file/internal/domain/models"
)

func TestParser_Parse(t *testing.T) {
	// Load test data
	validData, err := os.ReadFile("../../../test/samples/petstore.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	invalidJSON := []byte("{invalid:json}")

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid swagger",
			data:    validData,
			wantErr: false,
		},
		{
			name:    "invalid json",
			data:    invalidJSON,
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
	}

	parser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Parse(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("Parser.Parse() returned nil but expected a document")
			}
		})
	}
}

func TestParser_Validate(t *testing.T) {
	// Load test data
	validData, err := os.ReadFile("../../../test/samples/petstore.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	invalidData, err := os.ReadFile("../../../test/samples/invalid.json")
	if err != nil {
		t.Fatalf("Failed to read invalid test file: %v", err)
	}

	parser := New()

	// Parse valid document
	validDoc, err := parser.Parse(validData)
	if err != nil {
		t.Fatalf("Failed to parse valid data: %v", err)
	}

	// Parse invalid document
	invalidDoc, err := parser.Parse(invalidData)
	if err != nil {
		t.Fatalf("Failed to parse invalid data: %v", err)
	}

	// Create a document with missing paths
	missingPathsDoc := &models.SwaggerDoc{
		Swagger: "2.0",
		Info: models.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: map[string]models.PathItem{},
	}

	// Create an invalid path item with incomplete operation
	invalidPathDoc := &models.SwaggerDoc{
		Swagger: "2.0",
		Info: models.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: map[string]models.PathItem{
			"/test": {
				Get: &models.Operation{
					Responses: map[string]models.Response{},
				},
			},
		},
	}

	// Create an invalid path parameter
	invalidPathParamDoc := &models.SwaggerDoc{
		Swagger: "2.0",
		Info: models.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: map[string]models.PathItem{
			"/test/{id}": {
				Get: &models.Operation{
					Parameters: []models.Parameter{
						{
							Name:     "id",
							In:       "path",
							Required: false, // Path parameters must be required
						},
					},
					Responses: map[string]models.Response{
						"200": {
							Description: "OK",
						},
					},
				},
			},
		},
	}

	// Create a document with missing info fields
	missingInfoFieldsDoc := &models.SwaggerDoc{
		Swagger: "2.0",
		Info:    models.Info{}, // Empty info
		Paths: map[string]models.PathItem{
			"/test": {
				Get: &models.Operation{
					Responses: map[string]models.Response{
						"200": {
							Description: "OK",
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name    string
		doc     *models.SwaggerDoc
		wantErr bool
	}{
		{
			name:    "valid document",
			doc:     validDoc,
			wantErr: false,
		},
		{
			name:    "invalid document - missing swagger version",
			doc:     invalidDoc,
			wantErr: true,
		},
		{
			name:    "nil document",
			doc:     nil,
			wantErr: true,
		},
		{
			name:    "missing paths",
			doc:     missingPathsDoc,
			wantErr: true,
		},
		{
			name:    "invalid path - no responses",
			doc:     invalidPathDoc,
			wantErr: true,
		},
		{
			name:    "invalid path parameter - not required",
			doc:     invalidPathParamDoc,
			wantErr: true,
		},
		{
			name:    "missing info fields",
			doc:     missingInfoFieldsDoc,
			wantErr: false, // Should NOT error anymore with our fix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.Validate(tt.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParser_GetBaseURL(t *testing.T) {
	tests := []struct {
		name string
		doc  *models.SwaggerDoc
		want string
	}{
		{
			name: "swagger v2 with scheme, host and basePath",
			doc: &models.SwaggerDoc{
				Schemes:  []string{"https"},
				Host:     "api.example.com",
				BasePath: "/v1",
			},
			want: "https://api.example.com/v1",
		},
		{
			name: "swagger with host and basePath, no scheme",
			doc: &models.SwaggerDoc{
				Host:     "api.example.com",
				BasePath: "/v1",
			},
			want: "http://api.example.com/v1",
		},
		{
			name: "openapi v3 with servers",
			doc: &models.SwaggerDoc{
				Servers: []models.Server{
					{URL: "https://staging.example.com/api"},
					{URL: "https://api.example.com/v1"},
				},
			},
			want: "https://staging.example.com/api",
		},
		{
			name: "fallback to default",
			doc:  &models.SwaggerDoc{},
			want: "http://localhost",
		},
	}

	parser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.GetBaseURL(tt.doc); got != tt.want {
				t.Errorf("Parser.GetBaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_ExtractOperations(t *testing.T) {
	// Load test data
	data, err := os.ReadFile("../../../test/samples/petstore.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	parser := New()

	// Parse valid document
	doc, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("Failed to parse data: %v", err)
	}

	// Extract operations
	operations := parser.ExtractOperations(doc)

	// Expected results
	tests := []struct {
		name              string
		tag               string
		expectedCount     int
		containsPath      string
		containsMethod    string
		expectedContains  bool
	}{
		{
			name:              "pets tag operations",
			tag:               "pets",
			expectedCount:     5, // GET /pets, POST /pets, GET /pets/{petId}, PUT /pets/{petId}, DELETE /pets/{petId}
			containsPath:      "/pets",
			containsMethod:    "GET",
			expectedContains:  true,
		},
		{
			name:              "non-existent tag",
			tag:               "non-existent",
			expectedCount:     0,
			containsPath:      "",
			containsMethod:    "",
			expectedContains:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagOps, found := operations[tt.tag]
			
			// Check if tag exists and count matches
			if tt.expectedCount > 0 && !found {
				t.Errorf("Tag %s not found in operations", tt.tag)
				return
			}
			
			if found && len(tagOps) != tt.expectedCount {
				t.Errorf("Expected %d operations for tag %s, got %d", tt.expectedCount, tt.tag, len(tagOps))
			}
			
			// Check if specific operation exists
			if tt.expectedContains && found {
				hasOperation := false
				for _, op := range tagOps {
					if op.Path == tt.containsPath && op.Method == tt.containsMethod {
						hasOperation = true
						break
					}
				}
				if !hasOperation {
					t.Errorf("Expected to find %s %s operation in tag %s, but did not", 
						tt.containsMethod, tt.containsPath, tt.tag)
				}
			}
		})
	}
}
