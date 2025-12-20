package sidebar

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// RenderLinks generates navigation links for a list of nodes using the Visitor pattern.
func RenderLinks(nodes []ast.Node, parentPath path.Path, values map[string]any) []g.Node {
	tree := &ast.Tree{Nodes: nodes}
	state := NewLinkState()

	visitor := &SidebarVisitor{values: values}
	walker := walk.NewWalker(visitor)

	// Walk the tree and collect links
	_ = walker.Walk(tree, state)

	return state.Links()
}
