package render

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/components/inputs"
)

// PrimitiveRenderer renders primitive fields
type PrimitiveRenderer struct{}

func (r *PrimitiveRenderer) Render(node schema.Node, ctx Context) g.Node {
	fieldName := ctx.Path.String()

	return h.Div(
		h.Class("field"),
		h.ID("field-"+fieldName),
		renderLabel(node, fieldName),
		renderPrimitiveInput(node, fieldName, ctx.Value),
		renderHelp(node),
	)
}

func renderPrimitiveInput(node schema.Node, name string, value any) g.Node {
	field := tags.Field{
		Name:      name,
		InputType: node.InputType,
		Required:  node.Required,
		Min:       node.Min,
		Max:       node.Max,
		Step:      node.Step,
		Options:   node.Options,
	}

	switch node.InputType {
	case tags.TypeCheckbox:
		return inputs.Checkbox(field, value)
	case tags.TypeNumber:
		return inputs.Number(field, value)
	case tags.TypePassword:
		return inputs.Password(field, value)
	case tags.TypeRange:
		return inputs.Range(field, value)
	case tags.TypeDate:
		return inputs.Date(field, value)
	case tags.TypeTime:
		return inputs.Time(field, value)
	case tags.TypeRadio:
		return inputs.Radio(field, value)
	case tags.TypeSelect:
		return inputs.Select(field, value)
	default:
		return inputs.Text(field, value)
	}
}

func renderLabel(node schema.Node, name string) g.Node {
	var required g.Node
	if node.Required {
		required = h.Span(h.Class("field__label-required"), g.Text("*"))
	}

	return h.Label(
		h.For(name),
		h.Class("field__label"),
		g.Text(node.Name),
		required,
	)
}

func renderHelp(node schema.Node) g.Node {
	if node.Help == "" {
		return nil
	}
	return h.P(h.Class("field__help"), g.Text(node.Help))
}
