package sidebar

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/schema"
)

// renderStructLink generates a collapsible navigation link for a struct field.
func renderStructLink(node schema.Node, path schema.Path, values map[string]any) g.Node {
	children := RenderLinks(node.Children, path, values)

	return h.Li(
		h.Class("nav__item nav__item--collapsible"),
		h.Div(
			h.Class("nav__item-header"),
			g.Attr("onclick", "toggleNavItem(this)"),
			h.Span(h.Class("nav__chevron"), g.Text("▾")),
			h.A(
				h.Href("#section-"+path.String()),
				h.Class("nav__link nav__link--section"),
				g.Text(node.Name),
			),
		),
		h.Ul(
			h.Class("nav__sublist"),
			g.Group(children),
		),
	)
}
