package sidebar

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

type TreeState struct {
	nodes []g.Node
}

func (s *TreeState) Append(node g.Node) {
	s.nodes = append(s.nodes, node)
}

func (s *TreeState) Output() g.Node {
	return g.Group(s.nodes)
}

func RenderTree(nodes []ast.Node, currentFocus path.Path, values map[string]any) g.Node {
	tree := &ast.Tree{Nodes: nodes}
	state := &TreeState{}

	visitor := &TreeVisitor{
		currentFocus: currentFocus,
		values:       values,
	}

	walker := walk.NewWalker(visitor)
	_ = walker.Walk(tree, state)

	return h.Div(
		h.Class("sidebar-tree"),
		h.Div(
			h.Class("tree-root"),
			h.A(h.Class("nav__link"), h.Href("/"),
				g.Text("Config"),
			),
		),
		state.Output(),
	)
}
