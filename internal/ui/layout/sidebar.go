package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ui/layout/sidebar"
)

func Sidebar(s ast.Schema, values path.ValuesByPath, focus path.Path) g.Node {
	return h.Aside(
		h.Class("app__sidebar"),
		h.Div(
			h.Class("nav"),
			h.H3(h.Class("nav__title"), g.Text("On this page")),
			sidebar.RenderTree(s.Nodes, focus, values),
		),
	)
}
