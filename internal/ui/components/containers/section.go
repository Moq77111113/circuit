package containers

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Section(title string, children []g.Node, collapsible bool, collapsed bool) g.Node {
	classes := "section"
	if collapsed {
		classes += " collapsed"
	}

	var header g.Node
	if collapsible {
		header = h.Div(
			h.Class("section__header"),
			g.Attr("onclick", "toggleCollapse(this)"),
			h.Span(h.Class("section__chevron"), g.Text("▼")),
			h.H3(h.Class("section__title"), g.Text(title)),
		)
	} else {
		header = h.H3(h.Class("section__title"), g.Text(title))
	}

	body := h.Div(h.Class("section__body"), g.Group(children))

	return h.Section(
		h.ID("section-"+title),
		h.Class(classes),
		header,
		body,
	)
}
