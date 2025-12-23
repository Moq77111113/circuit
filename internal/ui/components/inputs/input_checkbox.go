package inputs

import (
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
