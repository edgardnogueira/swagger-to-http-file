package models

// Basic structure placeholders for the Swagger document model
// Will be expanded during implementation

type SwaggerDoc struct {
	Info       Info                  `json:"info"`
	BasePath   string                `json:"basePath"`
	Host       string                `json:"host"`
	Schemes    []string              `json:"schemes"`
	Paths      map[string]PathItem   `json:"paths"`
	Definitions map[string]SchemaObj `json:"definitions,omitempty"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
}

type Operation struct {
	OperationID string        `json:"operationId,omitempty"`
	Summary     string        `json:"summary,omitempty"`
	Description string        `json:"description,omitempty"`
	Parameters  []Parameter   `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
}

type Parameter struct {
	Name        string    `json:"name"`
	In          string    `json:"in"` // path, query, header, body, form
	Description string    `json:"description,omitempty"`
	Required    bool      `json:"required"`
	Schema      *SchemaObj `json:"schema,omitempty"`
	Type        string    `json:"type,omitempty"`
}

type Response struct {
	Description string    `json:"description"`
	Schema      *SchemaObj `json:"schema,omitempty"`
}

type SchemaObj struct {
	Ref         string              `json:"$ref,omitempty"`
	Type        string              `json:"type,omitempty"`
	Properties  map[string]SchemaObj `json:"properties,omitempty"`
	Items       *SchemaObj           `json:"items,omitempty"`
}
