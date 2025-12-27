package render

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// Render generates HTML for a tree of nodes with their values.
// Uses default rendering options.
func Render(nodes []ast.Node, values ast.ValuesByPath, focus path.Path) g.Node {
	return RenderWithOptions(nodes, values, focus, DefaultOptions())
}

// RenderWithOptions generates HTML for a tree of nodes with custom rendering options.
func RenderWithOptions(nodes []ast.Node, values ast.ValuesByPath, focus path.Path, opts Options) g.Node {
	tree := &ast.Tree{Nodes: nodes}

	visitor := &RenderVisitor{
		values:  values,
		options: opts,
		nodes:   []g.Node{},
	}
	walker := walk.NewWalker(visitor, walk.WithBasePath(focus))
	_ = walker.Walk(tree, nil)

	return g.Group(visitor.nodes)
}
