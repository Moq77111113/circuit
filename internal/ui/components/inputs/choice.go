package inputs

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Checkbox(field tags.Field, value any) g.Node {
	checked := value != nil && value.(bool)

	return h.Div(
		h.Class("toggle-switch"),
		h.Input(
			h.Type("radio"),
			h.Name(field.Name),
			h.Value("true"),
			h.ID(field.Name+"_on"),
			h.Class("toggle-switch__input"),
			g.If(checked, h.Checked()),
		),
		h.Input(
			h.Type("radio"),
			h.Name(field.Name),
			h.Value("false"),
			h.ID(field.Name+"_off"),
			h.Class("toggle-switch__input"),
			g.If(!checked, h.Checked()),
		),
		h.Label(
			h.For(field.Name+"_on"),
			h.Class("toggle-switch__label toggle-switch__label--on"),
			g.Text("On"),
		),
		h.Label(
			h.For(field.Name+"_off"),
			h.Class("toggle-switch__label toggle-switch__label--off"),
			g.Text("Off"),
		),
		h.Div(h.Class("toggle-switch__slider")),
	)
}

func Radio(field tags.Field, value any) g.Node {
	var options []g.Node
	currentVal := fmt.Sprintf("%v", value)

	for _, opt := range field.Options {
		attrs := []g.Node{
			h.Type("radio"),
			h.Name(field.Name),
			h.Value(opt.Value),
			h.ID(field.Name + "_" + opt.Value),
		}

		if field.Required {
			attrs = append(attrs, h.Required())
		}

		if value != nil && currentVal == opt.Value {
			attrs = append(attrs, h.Checked())
		}

		options = append(options, h.Div(
			h.Class("radio-option"),
			h.Input(attrs...),
			h.Label(
				h.For(field.Name+"_"+opt.Value),
				g.Text(opt.Label),
			),
		))
	}

	return h.Div(
		h.Class("radio-group"),
		g.Group(options),
	)
}

func Select(field tags.Field, value any) g.Node {
	var options []g.Node
	currentVal := fmt.Sprintf("%v", value)

	for _, opt := range field.Options {
		attrs := []g.Node{
			h.Value(opt.Value),
		}
		if value != nil && currentVal == opt.Value {
			attrs = append(attrs, h.Selected())
		}
		options = append(options, h.Option(
			g.Group(attrs),
			g.Text(opt.Label),
		))
	}

	attrs := []g.Node{
		h.Name(field.Name),
		h.ID(field.Name),
		h.Class("field__select"),
	}

	if field.Required {
		attrs = append(attrs, h.Required())
	}

	return h.Select(
		g.Group(attrs),
		g.Group(options),
	)
}
