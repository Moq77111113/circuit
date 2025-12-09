package inputs

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Slice(field tags.Field, value any) g.Node {
	items := reflectSliceValues(value)

	var nodes []g.Node

	if len(items) == 0 {
		nodes = append(nodes, h.Div(
			h.Class("slice__empty"),
			g.Text("No items"),
		))
	} else {
		for i, item := range items {
			nodes = append(nodes, renderSliceItem(field, i, item))
		}
	}

	nodes = append(nodes, h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("add:%s", field.Name)),
		h.Class("slice__add-btn"),
		g.Text("Add Item"),
	))

	return h.Div(h.Class("slice"), g.Group(nodes))
}

func renderSliceItem(field tags.Field, index int, value any) g.Node {
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

	var input g.Node
	switch field.InputType {
	case tags.TypeText:
		input = Text(indexedField, value)
	case tags.TypeNumber:
		input = Number(indexedField, value)
	case tags.TypeCheckbox:
		input = Checkbox(indexedField, value)
	default:
		input = Text(indexedField, value)
	}

	return h.Div(
		h.Class("slice__item"),
		input,
		h.Button(
			h.Type("submit"),
			h.Name("action"),
			h.Value(fmt.Sprintf("remove:%s:%d", field.Name, index)),
			h.Class("slice__remove-btn"),
			g.Text("Remove"),
		),
	)
}

