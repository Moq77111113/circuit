package node

import "github.com/moq77111113/circuit/internal/tags"

// UIMetadata holds rendering information separate from core AST.
type UIMetadata struct {
	InputType   tags.InputType
	Help        string
	Placeholder string
	Required    bool
	ReadOnly    bool
	Min         string
	Max         string
	Step        string
	Pattern     string
	MinLen      int
	MaxLen      int
	Options     []tags.Option
}

// Node represents a field in the config schema tree.
type Node struct {
	// Core AST
	ID       string
	Name     string
	Kind     NodeKind
	Children []Node
	Parent   *Node

	// Type info
	ValueType   ValueType // For primitives
	ElementKind NodeKind  // For slices

	// UI metadata (separated from core AST)
	UI *UIMetadata
}
