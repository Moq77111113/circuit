package collapsible

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ui/styles"
)

// Config configures a collapsible component.
type Config struct {
	ID        string
	Title     string
	Depth     int
	Count     int
	Collapsed bool
}

// Collapsible renders a collapsible container with BEM classes and CSS-based icons.
func Collapsible(cfg Config, children []g.Node) g.Node {
	classes := []string{styles.Collapsible, styles.DepthClass(cfg.Depth)}
	if cfg.Collapsed {
		classes = append(classes, styles.CollapsibleCollapsed)
	}

	attrs := []g.Node{h.Class(strings.Join(classes, " "))}
	if cfg.ID != "" {
		attrs = append(attrs, h.ID(cfg.ID))
	}
	attrs = append(attrs, Header(cfg.Title, cfg.Count), Body(children))

	return h.Div(attrs...)
}

// Header renders the collapsible header with title and optional count badge.
func Header(title string, count int) g.Node {
	children := []g.Node{
		h.Span(h.Class(styles.CollapsibleIcon + " " + styles.IconChevronDown)),
		h.Span(h.Class(styles.CollapsibleTitle), g.Text(title)),
	}

	if count > 0 {
		children = append(children, h.Span(
			h.Class(styles.CollapsibleCount),
			g.Textf("(%d)", count),
		))
	}

	return h.Div(
		h.Class(styles.CollapsibleHeader),
		g.Group(children),
	)
}

// Body renders the collapsible body containing child elements.
func Body(children []g.Node) g.Node {
	return h.Div(h.Class(styles.CollapsibleBody), g.Group(children))
}
