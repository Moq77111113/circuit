package render

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
	"github.com/moq77111113/circuit/internal/ui/components/collapsible"
	"github.com/moq77111113/circuit/internal/ui/components/containers"
)

// renderStructSliceItemWithFields renders a complete struct slice item with fields and remove button.
func (v *RenderVisitor) renderStructSliceItemWithFields(ctx *walk.VisitContext, node *ast.Node, index int, itemPath path.Path) g.Node {
	rc := ctx.Context.(*RenderContext)
	itemFields := v.renderStructSliceFields(ctx, node, index)

	itemValue := rc.Values[itemPath.String()]

	summary := containers.Extract(*node, itemValue, 2)
	summaryText := containers.Format(summary)

	var body []g.Node
	body = append(body, itemFields...)

	if !rc.ReadOnly {
		field, idx := parseItemPath(itemPath.String())
		removeButton := h.Button(
			h.Type("submit"),
			h.Name("action"),
			h.Value(fmt.Sprintf("remove:%s:%s", field, idx)),
			h.Class("button button--danger"),
			g.Text("Remove"),
		)
		body = append(body, removeButton)
	}

	cfg := collapsible.Config{
		ID:        fmt.Sprintf("slice-item-%s", itemPath.String()),
		Title:     fmt.Sprintf("#%d", index),
		Summary:   summaryText,
		Depth:     rc.ClampDepth(ctx.Depth + 1),
		Collapsed: true,
	}

	return collapsible.Collapsible(cfg, body)
}

// renderStructSliceFields renders all fields for a struct item in a slice.
func (v *RenderVisitor) renderStructSliceFields(ctx *walk.VisitContext, node *ast.Node, index int) []g.Node {
	itemPath := ctx.Path.Index(index)
	return v.renderFields(ctx, node.Children, itemPath)
}

// renderFields recursively renders fields (primitives and nested structs).
func (v *RenderVisitor) renderFields(ctx *walk.VisitContext, children []ast.Node, basePath path.Path) []g.Node {
	rc := ctx.Context.(*RenderContext)
	var fieldNodes []g.Node

	for i := range children {
		child := &children[i]
		childPath := basePath.Child(child.Name)

		switch child.Kind {
		case ast.KindPrimitive:
			value := rc.Values[childPath.String()]
			field := h.Div(
				h.Class("field"),
				h.ID("field-"+childPath.String()),
				renderLabel(child, childPath.String()),
				renderInput(child, childPath.String(), value, rc),
				renderHelp(child),
			)
			fieldNodes = append(fieldNodes, field)

		case ast.KindStruct:
			nestedFields := v.renderFields(ctx, child.Children, childPath)
			fieldNodes = append(fieldNodes, nestedFields...)
		}
	}

	return fieldNodes
}
