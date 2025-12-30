package layout

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ui/assets"
	"github.com/moq77111113/circuit/internal/ui/form"
	"github.com/moq77111113/circuit/internal/ui/layout/breadcrumb"
)

// Page renders a complete HTML page using a PageContext.
func Page(pc *PageContext) g.Node {
	formNode := form.Form(pc.RenderContext)

	title := pc.Title
	if title == "" {
		title = pc.Schema.Name + " Configuration"
	}

	bodyContent := []g.Node{renderMobileToggle(), renderMobileOverlay()}
	bodyContent = append(bodyContent, pc.TopContent...)

	mainContent := []g.Node{
		breadcrumb.RenderBreadcrumb(pc.Focus, pc.Schema.Nodes),
		h.Header(
			h.Class("header"),
			h.H1(h.Class("header__title"), g.Text(title)),
			h.P(h.Class("header__description"), g.Text("Configure your application settings below.")),
		),
	}

	if pc.ErrorMessage != "" {
		mainContent = append(mainContent, renderErrorBanner(pc.ErrorMessage))
	}

	mainContent = append(mainContent, formNode)

	if !pc.ReadOnly && len(pc.Actions) > 0 {
		mainContent = append(mainContent, renderActions(pc.Actions))
	}

	if pc.Brand {
		mainContent = append(mainContent, renderFooter())
	}

	bodyContent = append(bodyContent,
		h.Div(
			h.Class("app"),
			Sidebar(pc.RenderContext),
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
		Head:     renderHead(title),
		Body:     bodyContent,
	})
}
