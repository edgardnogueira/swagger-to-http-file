package models

// HttpRequest represents a single HTTP request in the .http file format
type HttpRequest struct {
	Name        string
	Method      string
	Path        string
	Headers     map[string]string
	Body        string
	Description string
	Vars        map[string]string
	Tag         string
}

// HttpFile represents a collection of HTTP requests to be saved in a .http file
type HttpFile struct {
	BaseURL     string
	GlobalVars  map[string]string
	Requests    []HttpRequest
	Tag         string
}
