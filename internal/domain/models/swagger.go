package models

// SwaggerDoc represents the top-level Swagger/OpenAPI document structure
type SwaggerDoc struct {
	Swagger     string                `json:"swagger,omitempty"`
	OpenAPI     string                `json:"openapi,omitempty"`
	Info        Info                  `json:"info"`
	BasePath    string                `json:"basePath,omitempty"`
	Host        string                `json:"host,omitempty"`
	Schemes     []string              `json:"schemes,omitempty"`
	Paths       map[string]PathItem   `json:"paths"`
	Definitions map[string]SchemaObj  `json:"definitions,omitempty"`
	Components  *Components           `json:"components,omitempty"`
	Servers     []Server              `json:"servers,omitempty"`
	Tags        []Tag                 `json:"tags,omitempty"`
}

// Info contains metadata about the API
type Info struct {
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	Version        string   `json:"version"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
}

// Contact information for the API
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// Server represents a server object in OpenAPI v3
type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

// ServerVariable is a variable for server URL template substitution
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

// Tag provides metadata about the API tags
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Components contains the reusable components in OpenAPI v3
type Components struct {
	Schemas    map[string]SchemaObj    `json:"schemas,omitempty"`
	Parameters map[string]Parameter    `json:"parameters,omitempty"`
	Responses  map[string]Response     `json:"responses,omitempty"`
	Examples   map[string]interface{}  `json:"examples,omitempty"`
	Headers    map[string]Header       `json:"headers,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

// Header represents a header parameter
type Header struct {
	Description string     `json:"description,omitempty"`
	Required    bool       `json:"required,omitempty"`
	Schema      *SchemaObj `json:"schema,omitempty"`
}

// SecurityScheme defines a security scheme that can be used by operations
type SecurityScheme struct {
	Type        string `json:"type"`                  // "apiKey", "http", "oauth2", "openIdConnect"
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`        // for apiKey
	In          string `json:"in,omitempty"`          // for apiKey: "query", "header", "cookie"
	Scheme      string `json:"scheme,omitempty"`      // for http: "basic", "bearer"
	BearerFormat string `json:"bearerFormat,omitempty"` // for http: "bearer"
}

// PathItem describes the operations available on a single path
type PathItem struct {
	Ref         string      `json:"$ref,omitempty"`
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Get         *Operation  `json:"get,omitempty"`
	Put         *Operation  `json:"put,omitempty"`
	Post        *Operation  `json:"post,omitempty"`
	Delete      *Operation  `json:"delete,omitempty"`
	Options     *Operation  `json:"options,omitempty"`
	Head        *Operation  `json:"head,omitempty"`
	Patch       *Operation  `json:"patch,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty"`
}

// Operation describes a single API operation on a path
type Operation struct {
	Tags        []string               `json:"tags,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Description string                 `json:"description,omitempty"`
	OperationID string                 `json:"operationId,omitempty"`
	Consumes    []string               `json:"consumes,omitempty"`
	Produces    []string               `json:"produces,omitempty"`
	Parameters  []Parameter            `json:"parameters,omitempty"`
	RequestBody *RequestBody           `json:"requestBody,omitempty"`
	Responses   map[string]Response    `json:"responses"`
	Security    []map[string][]string  `json:"security,omitempty"`
	Deprecated  bool                   `json:"deprecated,omitempty"`
}

// RequestBody represents a request body in OpenAPI v3
type RequestBody struct {
	Description string                  `json:"description,omitempty"`
	Content     map[string]MediaTypeObj `json:"content"`
	Required    bool                    `json:"required,omitempty"`
}

// MediaTypeObj represents a media type object in OpenAPI v3
type MediaTypeObj struct {
	Schema   *SchemaObj             `json:"schema,omitempty"`
	Examples map[string]interface{} `json:"examples,omitempty"`
}

// Parameter describes a single operation parameter
type Parameter struct {
	Name            string      `json:"name"`
	In              string      `json:"in"` // query, header, path, cookie, body
	Description     string      `json:"description,omitempty"`
	Required        bool        `json:"required,omitempty"`
	Schema          *SchemaObj  `json:"schema,omitempty"`
	Type            string      `json:"type,omitempty"`      // string, number, integer, boolean, array, object
	Format          string      `json:"format,omitempty"`
	Items           *SchemaObj  `json:"items,omitempty"`     // for array type
	Enum            []interface{} `json:"enum,omitempty"`
	Default         interface{} `json:"default,omitempty"`
	Example         interface{} `json:"example,omitempty"`
	Style           string      `json:"style,omitempty"`     // OpenAPI v3
	Explode         bool        `json:"explode,omitempty"`   // OpenAPI v3
	AllowReserved   bool        `json:"allowReserved,omitempty"` // OpenAPI v3
}

// Response describes a single response from an API Operation
type Response struct {
	Description string                 `json:"description"`
	Schema      *SchemaObj             `json:"schema,omitempty"`
	Headers     map[string]Header      `json:"headers,omitempty"`
	Content     map[string]MediaTypeObj `json:"content,omitempty"` // OpenAPI v3
}

// SchemaObj describes a schema for request/response bodies and parameters
type SchemaObj struct {
	Ref           string                `json:"$ref,omitempty"`
	Type          string                `json:"type,omitempty"`
	Format        string                `json:"format,omitempty"`
	Title         string                `json:"title,omitempty"`
	Description   string                `json:"description,omitempty"`
	Default       interface{}           `json:"default,omitempty"`
	MultipleOf    float64               `json:"multipleOf,omitempty"`
	Maximum       float64               `json:"maximum,omitempty"`
	Minimum       float64               `json:"minimum,omitempty"`
	MaxLength     int                   `json:"maxLength,omitempty"`
	MinLength     int                   `json:"minLength,omitempty"`
	Pattern       string                `json:"pattern,omitempty"`
	MaxItems      int                   `json:"maxItems,omitempty"`
	MinItems      int                   `json:"minItems,omitempty"`
	UniqueItems   bool                  `json:"uniqueItems,omitempty"`
	MaxProperties int                   `json:"maxProperties,omitempty"`
	MinProperties int                   `json:"minProperties,omitempty"`
	Required      []string              `json:"required,omitempty"`
	Enum          []interface{}         `json:"enum,omitempty"`
	Items         *SchemaObj            `json:"items,omitempty"`
	Properties    map[string]SchemaObj  `json:"properties,omitempty"`
	AdditionalProperties interface{}    `json:"additionalProperties,omitempty"`
	Example       interface{}           `json:"example,omitempty"`
}
