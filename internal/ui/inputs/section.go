package inputs

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Section(title string, children []g.Node) g.Node {
	return h.Section(
		h.ID("section-"+title),
		h.Class("section"),
		h.H3(h.Class("section__title"), g.Text(title)),
		h.Div(
			h.Class("section__content"),
			g.Group(children),
		),
	)
}
