package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Footer() g.Node {
	return h.Footer(
		h.Class("footer"),
		h.P(
			g.Text("Powered by "),
			h.A(
				h.Href("https://github.com/moq77111113/circuit"),
				h.Target("_blank"),
				h.Rel("noopener noreferrer"),
				h.Class("footer__link"),
				g.Text("Circuit"),
			),
		),
	)
}
