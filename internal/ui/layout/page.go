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
	"github.com/moq77111113/circuit/internal/validation"
)

type PageOptions struct {
	Title      string
	Brand      bool
	ReadOnly   bool
	TopContent []g.Node
}

type pageConfig struct {
	schema     ast.Schema
	values     ast.ValuesByPath
	title      string
	brand      bool
	focus      path.Path
	formNode   g.Node
	topContent []g.Node
}

func Page(s ast.Schema, values ast.ValuesByPath, focus path.Path, opts PageOptions) g.Node {
	formNode := form.Form(s, values, focus, opts.ReadOnly)
	return renderPage(pageConfig{s, values, opts.Title, opts.Brand, focus, formNode, opts.TopContent})
}

func PageWithErrors(s ast.Schema, values ast.ValuesByPath, focus path.Path, errors *validation.ValidationResult, opts PageOptions) g.Node {
	formNode := form.FormWithErrors(s, values, focus, errors, opts.ReadOnly)
	return renderPage(pageConfig{s, values, opts.Title, opts.Brand, focus, formNode, nil})
}

func renderPage(cfg pageConfig) g.Node {
	if cfg.title == "" {
		cfg.title = cfg.schema.Name + " Configuration"
	}

	bodyContent := []g.Node{renderMobileToggle(), renderMobileOverlay()}
	bodyContent = append(bodyContent, cfg.topContent...)

	mainContent := []g.Node{
		breadcrumb.RenderBreadcrumb(cfg.focus, cfg.schema.Nodes),
		h.Header(
			h.Class("header"),
			h.H1(h.Class("header__title"), g.Text(cfg.title)),
			h.P(h.Class("header__description"), g.Text("Configure your application settings below.")),
		),
		cfg.formNode,
	}

	if cfg.brand {
		mainContent = append(mainContent, renderFooter())
	}

	bodyContent = append(bodyContent,
		h.Div(
			h.Class("app"),
			Sidebar(cfg.schema, cfg.values, cfg.focus),
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
		Title:    cfg.title,
		Language: "en",
		Head:     renderHead(cfg.title),
		Body:     bodyContent,
	})
}
