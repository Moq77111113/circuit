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
	nodes []g.Node
}

// VisitPrimitive renders a primitive field.
func (v *RenderVisitor) VisitPrimitive(ctx *walk.VisitContext, node *ast.Node) error {
	rc := ctx.Context.(*RenderContext)
	value := rc.Values[ctx.Path.String()]

	var errorMessage string
	if rc.Errors != nil {
		errorMessage = rc.Errors.Get(ctx.Path)
	}

	field := h.Div(
		h.Class(styles.Field),
		h.ID("field-"+ctx.Path.String()),
		renderLabel(node, ctx.Path.String()),
		renderInput(node, ctx.Path.String(), value, rc),
		renderHelp(node),
		renderError(errorMessage),
	)

	v.nodes = append(v.nodes, field)
	return nil
}

// VisitStruct renders a struct node.
func (v *RenderVisitor) VisitStruct(ctx *walk.VisitContext, node *ast.Node) error {
	rc := ctx.Context.(*RenderContext)
	if ctx.Depth == 0 && rc.ShowCardsAtDepth0 {
		card := RenderStructCard(*node, ctx.Path, rc.Values)
		v.nodes = append(v.nodes, card)
		return walk.ErrSkipChildren
	}

	return nil
}

// VisitSlice renders a slice with collapsible container.
func (v *RenderVisitor) VisitSlice(ctx *walk.VisitContext, node *ast.Node) error {
	rc := ctx.Context.(*RenderContext)
	value := rc.Values[ctx.Path.String()]
	items := reflection.SliceValues(value)

	isCollapsed := rc.ShouldCollapse(ctx.Depth)

	var itemNodes []g.Node
	if len(items) == 0 {
		itemNodes = append(itemNodes, renderEmptyState())
	} else {
		for i, itemValue := range items {
			itemPath := ctx.Path.Index(i)

			if node.ElementKind == ast.KindPrimitive {
				itemNodes = append(itemNodes, renderPrimitiveSliceItem(node, i, itemValue, itemPath, rc))
			} else {
				itemNodes = append(itemNodes, v.renderStructSliceItemWithFields(ctx, node, i, itemPath))
			}
		}
	}

	itemNodes = append(itemNodes, renderAddButton(ctx.Path, rc.ReadOnly))

	cfg := collapsible.Config{
		ID:        "slice-" + ctx.Path.String(),
		Title:     ast.DisplayName(node),
		Depth:     rc.ClampDepth(ctx.Depth),
		Count:     len(items),
		Collapsed: isCollapsed,
	}
	container := collapsible.Collapsible(cfg, itemNodes)

	v.nodes = append(v.nodes, container)
	return nil
}
