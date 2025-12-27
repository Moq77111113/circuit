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

func Page(s ast.Schema, values path.ValuesByPath, title string, brand bool, focus path.Path, topContent ...g.Node) g.Node {
	if title == "" {
		title = s.Name + " Configuration"
	}

	bodyContent := []g.Node{
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
	}

	bodyContent = append(bodyContent, topContent...)

	mainContent := []g.Node{
		breadcrumb.RenderBreadcrumb(focus, s.Nodes),
		h.Header(
			h.Class("header"),
			h.H1(h.Class("header__title"), g.Text(title)),
			h.P(h.Class("header__description"), g.Text("Configure your application settings below.")),
		),
		form.Form(s, values, focus),
	}

	if brand {
		mainContent = append(mainContent, h.Footer(
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
		))
	}

	bodyContent = append(bodyContent,
		h.Div(
			h.Class("app"),
			Sidebar(s, values, focus),
			h.Main(
				h.Class("app__main"),
				h.Div(
					h.Class("app__container"),
					g.Group(mainContent),
				),
			),
		),
		h.Script(g.Raw(assets.DefaultJS)),
	)

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
				h.Link(h.Rel("icon"), h.Type("image/svg+xml"), h.Href("data:image/svg+xml,%3Csvg%20xmlns='http://www.w3.org/2000/svg'%20viewBox='0%200%20100%20100'%3E%3Cdefs%3E%3Cfilter%20id='glow'%3E%3CfeGaussianBlur%20stdDeviation='2'%20result='blur'/%3E%3CfeMerge%3E%3CfeMergeNode%20in='blur'/%3E%3CfeMergeNode%20in='SourceGraphic'/%3E%3C/feMerge%3E%3C/filter%3E%3C/defs%3E%3Crect%20width='100'%20height='100'%20fill='%23000000'/%3E%3Cpath%20d='M0%2050%20L10%2030%20L20%2055%20L30%2020%20L40%2060%20L50%2025%20L60%2055%20L70%2015%20L80%2050%20L90%2035%20L100%2050'%20stroke='%23FFFFFF'%20stroke-width='6'%20stroke-linecap='round'%20fill='none'%20filter='url(%23glow)'/%3E%3C/svg%3E")),
				h.TitleEl(g.Text(title)),
				h.StyleEl(g.Raw(assets.DefaultCSS)),
			),
			h.Body(bodyContent...),
		},
	})
}
