package ui

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/schema"
)

func Page(s schema.Schema, values map[string]any, title string, brand bool) g.Node {
	if title == "" {
		title = s.Name + " Configuration"
	}

	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []g.Node{
			h.Head(
				h.Meta(h.Charset("utf-8")),
				h.Meta(
					h.Name("viewport"),
					h.Content("width=device-width, initial-scale=1"),
				),
				FavIcons(),
				h.TitleEl(g.Text(title)),
				h.StyleEl(g.Raw(DefaultCSS)),
			),
			h.Body(
				h.Div(
					h.Class("app"),
					Sidebar(s),
					h.Main(
						h.Class("app__main"),
						h.Div(
							h.Class("app__container"),
							Header(title, "Configure your application settings below."),
							Form(s, values),
							g.If(brand, Footer()),
						),
					),
				),
				h.Script(g.Raw(DefaultJS)),
			),
		},
	})
}
