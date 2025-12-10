package containers

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/components/inputs"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func renderStructItem(field tags.Field, index int, value any, depth int) g.Node {
	var fields []g.Node
	for _, subfield := range field.Fields {
		indexedName := fmt.Sprintf("%s.%d.%s", field.Name, index, subfield.Name)
		indexedField := subfield
		indexedField.Name = indexedName

		subval := extractStructField(value, subfield.Name)

		if subfield.IsSlice {
			fields = append(fields, SliceWithDepth(indexedField, subval, depth+1))
			continue
		} else if subfield.InputType == tags.TypeSection {
			fields = append(fields, renderNestedSection(indexedField, subval, depth+1))
			continue
		}

		var input g.Node
		switch subfield.InputType {
		case tags.TypeText:
			input = inputs.Text(indexedField, subval)
		case tags.TypeNumber:
			input = inputs.Number(indexedField, subval)
		case tags.TypeCheckbox:
			input = inputs.Checkbox(indexedField, subval)
		case tags.TypePassword:
			input = inputs.Password(indexedField, subval)
		case tags.TypeSelect:
			input = inputs.Select(indexedField, subval)
		case tags.TypeRadio:
			input = inputs.Radio(indexedField, subval)
		case tags.TypeRange:
			input = inputs.Range(indexedField, subval)
		default:
			input = inputs.Text(indexedField, subval)
		}

		fields = append(fields, h.Div(
			h.Class("field"),
			h.Label(
				h.For(indexedName),
				h.Class("field__label"),
				g.Text(subfield.Name),
			),
			input,
		))
	}

	return h.Div(h.Class("slice__struct"), g.Group(fields))
}
