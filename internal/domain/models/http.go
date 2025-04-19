package models

// HTTPRequest represents a single HTTP request in the .http file format
type HTTPRequest struct {
	Name        string
	Method      string
	Path        string
	Headers     map[string]string
	Body        string
	Description string
	Vars        map[string]string
	Tag         string
}

// HTTPFile represents a collection of HTTP requests to be saved in a .http file
type HTTPFile struct {
	BaseURL    string
	GlobalVars map[string]string
	Requests   []HTTPRequest
	Tag        string
}
