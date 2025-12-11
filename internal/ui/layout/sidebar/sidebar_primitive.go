package sidebar

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/schema"
)

// renderPrimitiveLink generates a navigation link for a primitive field.
func renderPrimitiveLink(node schema.Node, path schema.Path) g.Node {
	return h.Li(
		h.Class("nav__item"),
		h.A(
			h.Href("#field-"+path.String()),
			h.Class("nav__link nav__link--field"),
			g.Text(node.Name),
		),
	)
}
