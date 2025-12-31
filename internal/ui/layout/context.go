package layout

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/ui/render"
)

// ActionButton represents an executable action in the UI.
type ActionButton struct {
	Name                string
	Label               string
	Description         string
	RequireConfirmation bool
}

// PageContext extends RenderContext with page-level metadata.
type PageContext struct {
	*render.RenderContext

	// Page metadata
	Title        string
	Brand        bool
	TopContent   []g.Node
	Actions      []ActionButton
	ErrorMessage string
}

// NewPageContext creates a PageContext from a RenderContext.
func NewPageContext(rc *render.RenderContext) *PageContext {
	title := rc.Schema.Name + " Configuration"
	return &PageContext{
		RenderContext: rc,
		Title:         title,
		Brand:         false,
	}
}
