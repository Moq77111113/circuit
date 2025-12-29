package sidebar

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
	"github.com/moq77111113/circuit/internal/ui/render"
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

// RenderTree renders a sidebar navigation tree using a RenderContext.
func RenderTree(rc *render.RenderContext) g.Node {
	tree := &ast.Tree{Nodes: rc.Schema.Nodes}
	state := &TreeState{}

	visitor := &TreeVisitor{}

	walker := walk.NewWalker(visitor)
	_ = walker.WalkWithContext(tree, state, rc)

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
