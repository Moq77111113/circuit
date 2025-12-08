package ui

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/tags"
)

func Sidebar(s schema.Schema) g.Node {
	return h.Aside(
		h.Class("app__sidebar"),
		h.Div(
			h.Class("nav"),
			h.H3(h.Class("nav__title"), g.Text("On this page")),
			h.Ul(
				h.Class("nav__list"),
				g.Group(renderSidebarLinks(s.Fields)),
			),
		),
	)
}

func renderSidebarLinks(fields []tags.Field) []g.Node {
	var links []g.Node
	for _, field := range fields {
		if field.InputType == tags.TypeSection {
			links = append(links, h.Li(
				h.Class("nav__item"),
				h.A(
					h.Href("#section-"+field.Name),
					h.Class("nav__link nav__link--section"),
					g.Text(field.Name),
				),
				h.Ul(
					h.Class("nav__sublist"),
					g.Group(renderSidebarLinks(field.Fields)),
				),
			))
		} else {
			links = append(links, h.Li(
				h.Class("nav__item"),
				h.A(
					h.Href("#field-"+field.Name),
					h.Class("nav__link nav__link--field"),
					g.Text(field.Name),
				),
			))
		}
	}
	return links
}
