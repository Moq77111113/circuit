package layout

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ui/assets"
	"github.com/moq77111113/circuit/internal/ui/form"
	"github.com/moq77111113/circuit/internal/ui/layout/breadcrumb"
)

func Page(s ast.Schema, values map[string]any, title string, brand bool, focus path.Path) g.Node {
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
				Links(),
				h.TitleEl(g.Text(title)),
				h.StyleEl(g.Raw(assets.DefaultCSS)),
			),
			h.Body(
				h.Button(
					h.Class("mobile-menu-toggle"),
					h.Type("button"),
					g.Attr("aria-label", "Toggle menu"),
					g.Attr("onclick", "toggleSidebar()"),
					g.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M9 6l6 6-6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>`),
				),
				h.Div(
					h.Class("mobile-overlay"),
					g.Attr("onclick", "closeSidebar()"),
				),
				h.Div(
					h.Class("app"),
					Sidebar(s, values, focus),
					h.Main(
						h.Class("app__main"),
						h.Div(
							h.Class("app__container"),
							breadcrumb.RenderBreadcrumb(focus, s.Nodes),
							Header(title, "Configure your application settings below."),
							form.Form(s, values, focus),
							g.If(brand, Footer()),
						),
					),
				),
				h.Script(g.Raw(assets.DefaultJS)),
			),
		},
	})
}
