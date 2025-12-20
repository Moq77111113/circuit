package sidebar

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// SidebarVisitor implements walk.Visitor for navigation link generation.
type SidebarVisitor struct {
	values map[string]any
}

// VisitPrimitive generates a navigation link for a primitive field.
func (v *SidebarVisitor) VisitPrimitive(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*LinkState)

	link := h.Li(
		h.Class("nav__item"),
		h.A(
			h.Href("#field-"+ctx.Path.String()),
			h.Class("nav__link nav__link--field"),
			g.Text(node.Name),
		),
	)

	state.Add(link)
	return nil
}

// VisitStruct generates a navigation section link for a struct.
func (v *SidebarVisitor) VisitStruct(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*LinkState)

	// Add section link for the struct
	sectionLink := h.Li(
		h.Class("nav__item"),
		h.A(
			h.Href("#section-"+ctx.Path.String()),
			h.Class("nav__link nav__link--section"),
			g.Text(node.Name),
		),
	)

	state.Add(sectionLink)

	// Walker will visit children automatically
	return nil
}
