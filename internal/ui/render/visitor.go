package render

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
	"github.com/moq77111113/circuit/internal/reflection"
	"github.com/moq77111113/circuit/internal/ui/components/containers"
)

// RenderVisitor implements walk.Visitor for HTML rendering.
type RenderVisitor struct {
	values map[string]any // All configuration values indexed by path
}

// VisitPrimitive renders a primitive field (input + label + help).
func (v *RenderVisitor) VisitPrimitive(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*RenderState)
	value := v.values[ctx.Path.String()]

	field := h.Div(
		h.Class("field"),
		h.ID("field-"+ctx.Path.String()),
		renderLabel(node, ctx.Path.String()),
		renderInput(node, ctx.Path.String(), value),
		renderHelp(node),
	)

	state.Append(field)
	return nil
}

// VisitStruct renders a struct node.
func (v *RenderVisitor) VisitStruct(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*RenderState)

	if ctx.Depth == 0 {
		card := RenderStructCard(*node, ctx.Path, v.values)
		state.Append(card)
	}

	return nil
}

// VisitSlice renders a slice with collapsible container.
func (v *RenderVisitor) VisitSlice(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*RenderState)
	value := v.values[ctx.Path.String()]
	items := reflection.SliceValues(value)

	isCollapsed := ctx.Depth >= 2
	header := containers.CollapsibleHeader(node.Name, len(items), isCollapsed, "")

	var itemNodes []g.Node
	if len(items) == 0 {
		itemNodes = append(itemNodes, renderEmptyState())
	} else {
		for i, itemValue := range items {
			itemPath := ctx.Path.Index(i)

			if node.ElementKind == ast.KindPrimitive {
				itemNodes = append(itemNodes, renderPrimitiveSliceItem(node, i, itemValue, itemPath))
			} else {
				itemNodes = append(itemNodes, v.renderStructSliceItemWithFields(ctx, node, i, itemPath))
			}
		}
	}

	// Add button at the end
	itemNodes = append(itemNodes, renderAddButton(ctx.Path))

	body := containers.CollapsibleBody(itemNodes)
	container := containers.CollapsibleContainer(ctx.Depth, ctx.Path.String(), header, body)

	state.Append(container)
	return nil
}
