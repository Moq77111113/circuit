package sidebar

import (
	"reflect"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/schema"
)

// renderSliceLink generates a collapsible navigation link for a slice field.
func renderSliceLink(node schema.Node, path schema.Path, values map[string]any) g.Node {
	value := getValueAtPath(values, path.String())
	itemLinks := renderSliceItemLinks(node, path, value)

	return h.Li(
		h.Class("nav__item nav__item--collapsible"),
		h.Div(
			h.Class("nav__item-header"),
			g.Attr("onclick", "toggleNavItem(this)"),
			h.Span(h.Class("nav__chevron"), g.Text("▾")),
			h.A(
				h.Href("#slice-"+path.String()),
				h.Class("nav__link nav__link--slice"),
				g.Text(node.Name),
			),
		),
		h.Ul(
			h.Class("nav__sublist"),
			g.Group(itemLinks),
		),
	)
}

// renderSliceItemLinks generates links for actual slice items (0, 1, 2...).
func renderSliceItemLinks(node schema.Node, path schema.Path, value any) []g.Node {
	if value == nil {
		return nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice {
		return nil
	}

	var links []g.Node
	for i := 0; i < rv.Len(); i++ {
		itemPath := path.Index(i)
		links = append(links, renderSliceItemLink(node, itemPath, i))
	}
	return links
}

// renderSliceItemLink generates a link for a single slice item.
func renderSliceItemLink(node schema.Node, path schema.Path, index int) g.Node {
	href := "#field-" + path.String()
	class := "nav__link nav__link--field"

	if node.ElementKind == schema.KindStruct {
		href = "#slice-item-" + path.String()
		class = "nav__link nav__link--slice-item"
	}

	return h.Li(
		h.Class("nav__item"),
		h.A(
			h.Href(href),
			h.Class(class),
			g.Textf("Item %d", index),
		),
	)
}
