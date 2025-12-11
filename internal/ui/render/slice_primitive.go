package render

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/schema"
)

// renderPrimitiveItem renders a primitive slice item
func (r *SliceRenderer) renderPrimitiveItem(node schema.Node, index int, value any, path schema.Path) g.Node {
	removeBtn := h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("remove:%s:%d", path.String(), index)),
		h.Class("slice__remove-button"),
		g.Text("Remove"),
	)

	itemNode := schema.Node{
		Name:      path.String(),
		Kind:      schema.KindPrimitive,
		ValueType: node.ValueType,
		InputType: node.InputType,
		Required:  node.Required,
		Min:       node.Min,
		Max:       node.Max,
		Step:      node.Step,
	}

	itemCtx := Context{
		Path:  path,
		Value: value,
		Depth: 0,
	}

	input := (&PrimitiveRenderer{}).Render(itemNode, itemCtx)

	return h.Div(
		h.Class("slice__item slice__item--primitive"),
		input,
		removeBtn,
	)
}
