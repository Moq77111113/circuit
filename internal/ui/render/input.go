package render

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/components/inputs"
)

// valueTypeToString converts ValueType enum to string
func valueTypeToString(vt ast.ValueType) string {
	switch vt {
	case ast.ValueString:
		return "string"
	case ast.ValueInt:
		return "int"
	case ast.ValueBool:
		return "bool"
	case ast.ValueFloat:
		return "float"
	default:
		return "string"
	}
}

// renderLabel creates a label element for a field
func renderLabel(node *ast.Node, fieldName string) g.Node {
	return h.Label(
		h.For(fieldName),
		h.Class("field__label"),
		g.Text(node.Name),
	)
}

// renderInput creates an input element based on the node's InputType
func renderInput(node *ast.Node, fieldName string, value any) g.Node {
	// Create a tags.Field for compatibility with existing inputs components
	field := tags.Field{
		Name:      fieldName,
		Type:      valueTypeToString(node.ValueType),
		InputType: node.UI.InputType,
		Help:      node.UI.Help,
		Required:  node.UI.Required,
		Min:       node.UI.Min,
		Max:       node.UI.Max,
		Step:      node.UI.Step,
		Options:   node.UI.Options,
	}

	switch node.UI.InputType {
	case tags.TypeText, tags.TypeEmail, tags.TypeUrl, tags.TypePassword:
		return inputs.Text(field, value)
	case tags.TypeNumber:
		return inputs.Number(field, value)
	case tags.TypeCheckbox:
		return inputs.Checkbox(field, value)
	case tags.TypeDate:
		return inputs.Date(field, value)
	case tags.TypeTime:
		return inputs.Time(field, value)
	case tags.TypeRange:
		return inputs.Range(field, value)
	case tags.TypeSelect:
		return inputs.Select(field, value)
	case tags.TypeRadio:
		return inputs.Radio(field, value)
	default:
		return inputs.Text(field, value)
	}
}

// renderHelp creates a help text element if help is provided
func renderHelp(node *ast.Node) g.Node {
	if node.UI == nil || node.UI.Help == "" {
		return nil
	}
	return h.Span(
		h.Class("field__help"),
		g.Text(node.UI.Help),
	)
}
