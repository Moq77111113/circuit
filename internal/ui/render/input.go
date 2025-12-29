package render

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/components/inputs"
	"github.com/moq77111113/circuit/internal/ui/styles"
)

// renderLabel creates a label element for a field
func renderLabel(node *ast.Node, fieldName string) g.Node {
	return h.Label(
		h.For(fieldName),
		h.Class(styles.FieldLabel),
		g.Text(node.Name),
	)
}

// renderInput creates an input element based on the node's InputType
func renderInput(node *ast.Node, fieldName string, value any, rc *RenderContext) g.Node {
	field := tags.Field{
		Name:      fieldName,
		Type:      valueTypeToString(node.ValueType),
		InputType: node.UI.InputType,
		Help:      node.UI.Help,
		Required:  node.UI.Required,
		ReadOnly:  node.UI.ReadOnly || rc.ReadOnly,
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
		h.Class(styles.FieldHelp),
		g.Text(node.UI.Help),
	)
}

// renderError creates an error message element if message is provided
func renderError(message string) g.Node {
	if message == "" {
		return nil
	}
	return h.Span(
		h.Class(styles.FieldError),
		g.Text(message),
	)
}
