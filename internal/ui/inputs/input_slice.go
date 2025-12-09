package inputs

import (
	"fmt"
	"reflect"

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
	var content g.Node

	if len(field.Fields) > 0 {
		content = renderStructItem(field, index, value)
	} else {
		content = renderPrimitiveItem(field, index, value)
	}

	return h.Div(
		h.Class("slice__item"),
		content,
		h.Button(
			h.Type("submit"),
			h.Name("action"),
			h.Value(fmt.Sprintf("remove:%s:%d", field.Name, index)),
			h.Class("slice__remove-button"),
			g.Text("Remove"),
		),
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
		return Text(indexedField, value)
	case tags.TypeNumber:
		return Number(indexedField, value)
	case tags.TypeCheckbox:
		return Checkbox(indexedField, value)
	default:
		return Text(indexedField, value)
	}
}

func renderStructItem(field tags.Field, index int, value any) g.Node {
	var fields []g.Node
	for _, subfield := range field.Fields {
		indexedName := fmt.Sprintf("%s.%d.%s", field.Name, index, subfield.Name)
		indexedField := subfield
		indexedField.Name = indexedName

		subval := extractStructField(value, subfield.Name)

		var input g.Node
		switch subfield.InputType {
		case tags.TypeText:
			input = Text(indexedField, subval)
		case tags.TypeNumber:
			input = Number(indexedField, subval)
		case tags.TypeCheckbox:
			input = Checkbox(indexedField, subval)
		default:
			input = Text(indexedField, subval)
		}

		fields = append(fields, h.Div(
			h.Class("slice__struct-field"),
			h.Label(
				h.For(indexedName),
				h.Class("slice__struct-label"),
				g.Text(subfield.Name),
			),
			input,
		))
	}

	return h.Div(h.Class("slice__struct"), g.Group(fields))
}

func extractStructField(v any, fieldName string) any {
	if v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil
	}

	fv := rv.FieldByName(fieldName)
	if !fv.IsValid() {
		return nil
	}

	return fv.Interface()
}
