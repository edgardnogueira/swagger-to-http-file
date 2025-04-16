package models

// OperationInfo contains information about an API operation with its path and method
type OperationInfo struct {
	Path       string
	Method     string
	Operation  *Operation
	Parameters []Parameter
}
