package validation

import "github.com/moq77111113/circuit/internal/ast/path"

// ValidationError represents a single field validation error.
type ValidationError struct {
	Path    path.Path
	Field   string
	Message string
}

// ValidationResult holds the outcome of validating a form submission.
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// Has returns true if there is an error for the given path.
func (r *ValidationResult) Has(p path.Path) bool {
	for _, err := range r.Errors {
		if err.Path.String() == p.String() {
			return true
		}
	}
	return false
}

// Get returns the error message for the given path, or empty string if none.
func (r *ValidationResult) Get(p path.Path) string {
	for _, err := range r.Errors {
		if err.Path.String() == p.String() {
			return err.Message
		}
	}
	return ""
}

// FirstError returns the first validation error, or nil if there are none.
func (r *ValidationResult) FirstError() *ValidationError {
	if len(r.Errors) == 0 {
		return nil
	}
	return &r.Errors[0]
}
