package containers

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/components/inputs"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func renderNestedSection(field tags.Field, value any, depth int) g.Node {
	var children []g.Node
	for _, subfield := range field.Fields {
		nestedName := fmt.Sprintf("%s.%s", field.Name, subfield.Name)
		nestedField := subfield
		nestedField.Name = nestedName

		subval := extractStructField(value, subfield.Name)

		var input g.Node

		// Recursively handle slices and nested structs
		if subfield.IsSlice {
			input = SliceWithDepth(nestedField, subval, depth+1)
		} else if subfield.InputType == tags.TypeSection {
			input = renderNestedSection(nestedField, subval, depth+1)
		} else {
			switch subfield.InputType {
			case tags.TypeText:
				input = inputs.Text(nestedField, subval)
			case tags.TypeNumber:
				input = inputs.Number(nestedField, subval)
			case tags.TypeCheckbox:
				input = inputs.Checkbox(nestedField, subval)
			case tags.TypePassword:
				input = inputs.Password(nestedField, subval)
			case tags.TypeSelect:
				input = inputs.Select(nestedField, subval)
			case tags.TypeRadio:
				input = inputs.Radio(nestedField, subval)
			case tags.TypeRange:
				input = inputs.Range(nestedField, subval)
			default:
				input = inputs.Text(nestedField, subval)
			}
		}

		children = append(children, h.Div(
			h.Class("field"),
			h.Label(
				h.For(nestedName),
				h.Class("field__label"),
				g.Text(subfield.Name),
			),
			input,
		))
	}

	return Section(field.Name, children, true, false)
}
