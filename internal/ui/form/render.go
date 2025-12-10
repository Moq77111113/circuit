package form

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/components/containers"
	"github.com/moq77111113/circuit/internal/ui/components/inputs"
)

var renderers = map[tags.InputType]func(tags.Field, any) g.Node{
	tags.TypeCheckbox: inputs.Checkbox,
	tags.TypeNumber:   inputs.Number,
	tags.TypePassword: inputs.Password,
	tags.TypeRange:    inputs.Range,
	tags.TypeDate:     inputs.Date,
	tags.TypeTime:     inputs.Time,
	tags.TypeRadio:    inputs.Radio,
	tags.TypeSelect:   inputs.Select,
	tags.TypeText:     inputs.Text,
}

func renderField(field tags.Field, values map[string]any) g.Node {
	var value any
	if values != nil {
		value = values[field.Name]
	}

	if field.InputType == tags.TypeSection || field.IsSlice {
		return renderInput(field, value)
	}

	return h.Div(
		h.Class("field"),
		h.ID("field-"+field.Name),
		renderLabel(field),
		renderInput(field, value),
		renderHelp(field),
	)
}

func renderLabel(field tags.Field) g.Node {
	var required g.Node
	if field.Required {
		required = h.Span(h.Class("field__label-required"), g.Text("*"))
	}

	return h.Label(
		h.For(field.Name),
		h.Class("field__label"),
		g.Text(field.Name),
		required,
	)
}

func renderInput(field tags.Field, value any) g.Node {
	if field.InputType == tags.TypeSection {
		return renderSection(field, value)
	}
	if field.IsSlice {
		return containers.Slice(field, value)
	}
	if renderer, ok := renderers[field.InputType]; ok {
		return renderer(field, value)
	}
	return inputs.Text(field, value)
}

func renderHelp(field tags.Field) g.Node {
	if field.Help == "" {
		return nil
	}

	return h.P(
		h.Class("field__help"),
		g.Text(field.Help),
	)
}
