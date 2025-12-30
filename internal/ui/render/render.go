package render

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// Render generates HTML using a RenderContext and filtered nodes.
func Render(nodes []ast.Node, rc *RenderContext) g.Node {
	tree := &ast.Tree{Nodes: nodes}

	visitor := &RenderVisitor{
		nodes: []g.Node{},
	}
	walker := walk.NewWalker(visitor, walk.WithBasePath(rc.Focus))
	_ = walker.WalkWithContext(tree, nil, rc)

	return g.Group(visitor.nodes)
}
