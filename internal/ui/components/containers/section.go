package containers

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Section(title string, children []g.Node, collapsible bool, collapsed bool) g.Node {
	if collapsible {
		header := sectionCollapsibleHeader(title)
		body := sectionCollapsibleBody(children)

		containerClass := "container container--depth-0"
		if collapsed {
			containerClass += " collapsed"
		}

		return h.Section(
			h.ID("section-"+title),
			h.Class(containerClass),
			header,
			body,
		)
	}

	return h.Section(
		h.ID("section-"+title),
		h.Class("section"),
		h.H3(h.Class("section__title"), g.Text(title)),
		g.Group(children),
	)
}

func sectionCollapsibleHeader(title string) g.Node {
	displayTitle := title
	if strings.Contains(title, ".") {
		parts := strings.Split(title, ".")
		displayTitle = parts[len(parts)-1]
	}

	return h.Div(
		h.Class("container__header"),
		g.Attr("onclick", "toggleCollapse(this)"),
		h.Span(h.Class("container__chevron"), g.Text("▼")),
		h.Span(h.Class("container__title"), g.Text(displayTitle)),
	)
}

func sectionCollapsibleBody(children []g.Node) g.Node {
	return h.Div(h.Class("container__body"), g.Group(children))
}
