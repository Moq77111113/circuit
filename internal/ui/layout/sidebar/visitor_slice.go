package sidebar

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
	"github.com/moq77111113/circuit/internal/reflection"
)

// VisitSlice generates navigation links for a slice and its items.
func (v *SidebarVisitor) VisitSlice(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*LinkState)
	value := v.values[ctx.Path.String()]

	// Add link for the slice container
	containerLink := h.Li(
		h.Class("nav__item nav__item--collapsible"),
		h.A(
			h.Href("#slice-"+ctx.Path.String()),
			h.Class("nav__link nav__link--slice"),
			g.Attr("onclick", "toggleNavItem(this)"),
			h.Span(h.Class("nav__chevron"), g.Text("▾")),
			g.Text(node.Name),
		),
	)

	state.Add(containerLink)

	if value == nil {
		// No items, no item links
		return nil
	}

	items := reflection.SliceValues(value)
	if len(items) == 0 {
		return nil
	}

	// Generate a link for each slice item
	for i := range items {
		itemPath := ctx.Path.Index(i)
		label := fmt.Sprintf("Item %d", i)

		itemLink := h.Li(
			h.Class("nav__item nav__item--nested"),
			h.A(
				h.Href("#slice-item-"+itemPath.String()),
				h.Class("nav__link nav__link--item"),
				g.Text(label),
			),
		)

		state.Add(itemLink)
	}

	return nil
}
