package containers

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// CollapsibleHeader renders the header for a collapsible section or slice.
func CollapsibleHeader(title string, count int, collapsed bool) g.Node {
	return h.Div(
		h.Class("slice__header"),
		g.Attr("onclick", "toggleCollapse(this)"),
		h.Span(h.Class("slice__chevron"), g.Text("▼")),
		h.Span(h.Class("slice__title"), g.Text(title)),
		h.Span(h.Class("slice__count"), g.Textf("(%d)", count)),
	)
}

// CollapsibleBody renders the body of a collapsible section or slice.
func CollapsibleBody(children []g.Node) g.Node {
	return h.Div(h.Class("slice__body"), g.Group(children))
}

// CollapsibleContainer renders the container for a collapsible section or slice.
func CollapsibleContainer(depth int, children ...g.Node) g.Node {
	containerClass := fmt.Sprintf("slice %s", DepthClass(depth))
	if IsCollapsed(depth) {
		containerClass += " collapsed"
	}
	return h.Div(h.Class(containerClass), g.Group(children))
}

// SectionHeader renders the header for a section, optionally collapsible.
func SectionHeader(title string, collapsible bool) g.Node {
	if collapsible {
		return h.Div(
			h.Class("section__header"),
			g.Attr("onclick", "toggleCollapse(this)"),
			h.Span(h.Class("section__chevron"), g.Text("▼")),
			h.H3(h.Class("section__title"), g.Text(title)),
		)
	}
	return h.H3(h.Class("section__title"), g.Text(title))
}
