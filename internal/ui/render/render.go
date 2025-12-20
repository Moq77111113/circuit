package render

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// Render generates HTML for a tree of nodes with their values.
// This is the main public API for the render package.
func Render(nodes []ast.Node, values map[string]any) g.Node {
	tree := &ast.Tree{Nodes: nodes}
	state := NewRenderState()

	visitor := &RenderVisitor{values: values}
	walker := walk.NewWalker(visitor)
	_ = walker.Walk(tree, state)

	return g.Group(state.Output())
}
