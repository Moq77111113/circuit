package walk

import (
	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/path"
)

// Visitor defines the interface for traversing an AST.
type Visitor interface {
	VisitPrimitive(ctx *VisitContext, n *node.Node) error
	VisitStruct(ctx *VisitContext, n *node.Node) error
	VisitSlice(ctx *VisitContext, n *node.Node) error
}

// VisitContext holds state during tree traversal.
type VisitContext struct {
	Tree   *node.Tree
	Path   path.Path
	Depth  int
	Parent *node.Node
	Index  int // -1 if not in slice

	State any // Visitor-specific state
}

// NewContext creates a new visit context.
func NewContext(tree *node.Tree, state any) *VisitContext {
	return &VisitContext{
		Tree:  tree,
		Path:  path.Path{},
		Depth: 0,
		Index: -1,
		State: state,
	}
}
