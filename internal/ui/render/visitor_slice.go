package render

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// renderStructSliceItemWithFields renders a complete struct slice item with fields and remove button.
func (v *RenderVisitor) renderStructSliceItemWithFields(ctx *walk.VisitContext, node *ast.Node, index int, itemPath path.Path) g.Node {
	itemHeader := h.Div(
		h.Class("slice-item__header"),
		g.Textf("#%d", index),
	)
	itemFields := v.renderStructSliceFields(ctx, node, index)

	field, idx := parseItemPath(itemPath.String())
	removeBtn := h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("remove:%s:%s", field, idx)),
		h.Class("btn btn--remove"),
		g.Text("Remove"),
	)

	return h.Div(
		h.Class("slice-item slice-item--struct"),
		itemHeader,
		g.Group(itemFields),
		removeBtn,
	)
}

// renderStructSliceFields renders all fields for a struct item in a slice.
func (v *RenderVisitor) renderStructSliceFields(ctx *walk.VisitContext, node *ast.Node, index int) []g.Node {
	itemPath := ctx.Path.Index(index)
	return v.renderFields(node.Children, itemPath)
}

// renderFields recursively renders fields (primitives and nested structs).
func (v *RenderVisitor) renderFields(children []ast.Node, basePath path.Path) []g.Node {
	var fieldNodes []g.Node

	for i := range children {
		child := &children[i]
		childPath := basePath.Child(child.Name)

		switch child.Kind {
		case ast.KindPrimitive:
			value := v.values[childPath.String()]
			field := h.Div(
				h.Class("field"),
				h.ID("field-"+childPath.String()),
				renderLabel(child, childPath.String()),
				renderInput(child, childPath.String(), value),
				renderHelp(child),
			)
			fieldNodes = append(fieldNodes, field)

		case ast.KindStruct:
			nestedFields := v.renderFields(child.Children, childPath)
			fieldNodes = append(fieldNodes, nestedFields...)
		}
	}

	return fieldNodes
}
