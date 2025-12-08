package schema

import "github.com/moq77111113/circuit/internal/tags"

// Schema represents the complete structure of a config.
type Schema struct {
	Name   string
	Fields []tags.Field
}
