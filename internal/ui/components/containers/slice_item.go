package containers

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/components/inputs"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func RenderSliceItem(field tags.Field, index int, value any, depth int) g.Node {
	isStruct := len(field.Fields) > 0
	removeBtn := h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("remove:%s:%d", field.Name, index)),
		h.Class("slice__remove-button"),
		g.Text("Remove"),
	)

	if isStruct {
		return h.Div(
			h.Class("slice__item slice__item--struct"),
			h.Div(
				h.Class("slice__item-header"),
				h.Span(h.Class("slice__item-title"), g.Text(fmt.Sprintf("#%d", index+1))),
				removeBtn,
			),
			h.Div(
				h.Class("slice__item-body"),
				renderStructItem(field, index, value, depth),
			),
		)
	}

	return h.Div(
		h.Class("slice__item slice__item--primitive"),
		renderPrimitiveItem(field, index, value),
		removeBtn,
	)
}

func renderPrimitiveItem(field tags.Field, index int, value any) g.Node {
	indexedName := fmt.Sprintf("%s.%d", field.Name, index)
	indexedField := tags.Field{
		Name:      indexedName,
		Type:      field.ElementType,
		InputType: field.InputType,
		Required:  field.Required,
		Min:       field.Min,
		Max:       field.Max,
		Step:      field.Step,
	}

	switch field.InputType {
	case tags.TypeText:
		return inputs.Text(indexedField, value)
	case tags.TypeNumber:
		return inputs.Number(indexedField, value)
	case tags.TypeCheckbox:
		return inputs.Checkbox(indexedField, value)
	default:
		return inputs.Text(indexedField, value)
	}
}
