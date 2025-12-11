package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/ui/layout/sidebar"
)

func Sidebar(s schema.Schema, values map[string]any) g.Node {
	return h.Aside(
		h.Class("app__sidebar"),
		h.Div(
			h.Class("nav"),
			h.H3(h.Class("nav__title"), g.Text("On this page")),
			h.Ul(
				h.Class("nav__list"),
				g.Group(sidebar.RenderLinks(s.Nodes, schema.Path{}, values)),
			),
		),
	)
}
