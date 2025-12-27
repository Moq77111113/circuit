package render

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
	"github.com/moq77111113/circuit/internal/reflection"
	"github.com/moq77111113/circuit/internal/ui/components/collapsible"
	"github.com/moq77111113/circuit/internal/ui/styles"
)

// RenderVisitor implements walk.Visitor for HTML rendering.
type RenderVisitor struct {
	values  ast.ValuesByPath
	options Options
	nodes   []g.Node
}

// VisitPrimitive renders a primitive field.
func (v *RenderVisitor) VisitPrimitive(ctx *walk.VisitContext, node *ast.Node) error {
	value := v.values[ctx.Path.String()]

	field := h.Div(
		h.Class(styles.Field),
		h.ID("field-"+ctx.Path.String()),
		renderLabel(node, ctx.Path.String()),
		renderInput(node, ctx.Path.String(), value),
		renderHelp(node),
	)

	v.nodes = append(v.nodes, field)
	return nil
}

// VisitStruct renders a struct node.
func (v *RenderVisitor) VisitStruct(ctx *walk.VisitContext, node *ast.Node) error {
	if ctx.Depth == 0 && v.options.ShowCardsAtDepth0 {
		card := RenderStructCard(*node, ctx.Path, v.values)
		v.nodes = append(v.nodes, card)
	}

	return nil
}

// VisitSlice renders a slice with collapsible container.
func (v *RenderVisitor) VisitSlice(ctx *walk.VisitContext, node *ast.Node) error {
	value := v.values[ctx.Path.String()]
	items := reflection.SliceValues(value)

	isCollapsed := v.options.ShouldCollapse(ctx.Depth)

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

	itemNodes = append(itemNodes, renderAddButton(ctx.Path))

	cfg := collapsible.Config{
		ID:        "slice-" + ctx.Path.String(),
		Title:     ast.DisplayName(node),
		Depth:     v.options.ClampDepth(ctx.Depth),
		Count:     len(items),
		Collapsed: isCollapsed,
	}
	container := collapsible.Collapsible(cfg, itemNodes)

	v.nodes = append(v.nodes, container)
	return nil
}
