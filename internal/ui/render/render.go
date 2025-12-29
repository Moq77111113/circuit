package render

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// Render generates HTML using a RenderContext.
func Render(rc *RenderContext) g.Node {
	tree := &ast.Tree{Nodes: rc.Schema.Nodes}

	visitor := &RenderVisitor{
		nodes: []g.Node{},
	}
	walker := walk.NewWalker(visitor, walk.WithBasePath(rc.Focus))
	_ = walker.WalkWithContext(tree, nil, rc)

	return g.Group(visitor.nodes)
}
