package render

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// Render generates HTML for a tree of nodes with their values.
func Render(nodes []ast.Node, values map[string]any, focus path.Path) g.Node {
	tree := &ast.Tree{Nodes: nodes}
	state := NewRenderState()

	visitor := &RenderVisitor{values: values}
	walker := walk.NewWalker(visitor, walk.WithBasePath(focus))
	_ = walker.Walk(tree, state)

	return g.Group(state.Output())
}
