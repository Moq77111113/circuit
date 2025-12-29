package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ui/layout/sidebar"
	"github.com/moq77111113/circuit/internal/ui/render"
)

// Sidebar renders the sidebar navigation using a RenderContext.
func Sidebar(rc *render.RenderContext) g.Node {
	return h.Aside(
		h.Class("app__sidebar"),
		h.Div(
			h.Class("nav"),
			h.H3(h.Class("nav__title"), g.Text("On this page")),
			sidebar.RenderTree(rc),
		),
	)
}
