package layout

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ui/assets"
	"github.com/moq77111113/circuit/internal/ui/form"
	"github.com/moq77111113/circuit/internal/ui/layout/breadcrumb"
	"github.com/moq77111113/circuit/internal/ui/styles"
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
		breadcrumb.RenderBreadcrumb(pc.Focus, pc.Schema.Nodes, pc.HTTPBasePath),
		renderHeader(title, pc.Actions, pc.ReadOnly),
	}

	if pc.ErrorMessage != "" {
		mainContent = append(mainContent, renderErrorBanner(pc.ErrorMessage))
	}

	mainContent = append(mainContent, formNode)

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

func renderHeader(title string, actions []ActionButton, readOnly bool) g.Node {
	headerContent := []g.Node{
		h.Div(
			h.Class("header__content"),
			h.H1(h.Class("header__title"), g.Text(title)),
			h.P(h.Class("header__description"), g.Text("Configure your application settings below.")),
		),
	}

	if !readOnly && len(actions) > 0 {
		headerContent = append(headerContent, renderActionsDropdown(actions))
	}

	return h.Header(h.Class("header"), g.Group(headerContent))
}

func renderActionsDropdown(actions []ActionButton) g.Node {
	items := make([]g.Node, len(actions))
	for i, action := range actions {
		items[i] = renderActionItem(action)
	}

	return h.Div(
		h.Class(styles.ActionsDropdown),
		h.Button(
			h.Class(styles.ActionsButton),
			h.Type("button"),
			g.Attr("onclick", "toggleActionsDropdown()"),
			g.Text("Actions"),
			h.Span(h.Class(styles.ActionsIcon), g.Text("⋮")),
		),
		h.Div(
			h.Class(styles.ActionsMenu),
			h.ID("actions-menu"),
			g.Group(items),
		),
	)
}

func renderActionItem(action ActionButton) g.Node {
	buttonAttrs := []g.Node{
		h.Type("submit"),
		h.Class(styles.ActionsMenuItem),
	}

	if action.RequireConfirmation {
		onClick := "return confirmAction('" + action.Label + "', '" + action.Description + "')"
		buttonAttrs = append(buttonAttrs, g.Attr("onclick", onClick))
	}

	itemContent := []g.Node{
		h.Div(h.Class(styles.ActionsMenuItemLabel), g.Text(action.Label)),
	}

	if action.Description != "" {
		itemContent = append(itemContent, h.Div(h.Class(styles.ActionsMenuItemDesc), g.Text(action.Description)))
	}

	return h.Form(
		h.Method("post"),
		h.Class(styles.ActionsMenuItemForm),
		h.Input(h.Type("hidden"), h.Name("action"), h.Value("execute:"+action.Name)),
		h.Button(g.Group(buttonAttrs), g.Group(itemContent)),
	)
}

func renderErrorBanner(msg string) g.Node {
	return h.Div(
		h.Class(styles.ErrorBanner),
		h.P(g.Text("Error: "+msg)),
	)
}
