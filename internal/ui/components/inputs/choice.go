package inputs

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/styles"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func appendFieldAttrs(attrs []g.Node, field tags.Field) []g.Node {
	if field.Required {
		attrs = append(attrs, h.Required())
	}
	if field.ReadOnly {
		attrs = append(attrs, h.Disabled())
	}
	return attrs
}

func Checkbox(field tags.Field, value any) g.Node {
	checked := value != nil && value.(bool)

	onAttrs := []g.Node{
		h.Type("radio"),
		h.Name(field.Name),
		h.Value("true"),
		h.ID(field.Name + "_on"),
		h.Class(styles.ToggleSwitchInput),
		g.If(checked, h.Checked()),
	}
	offAttrs := []g.Node{
		h.Type("radio"),
		h.Name(field.Name),
		h.Value("false"),
		h.ID(field.Name + "_off"),
		h.Class(styles.ToggleSwitchInput),
		g.If(!checked, h.Checked()),
	}

	onAttrs = appendFieldAttrs(onAttrs, field)
	offAttrs = appendFieldAttrs(offAttrs, field)

	return h.Div(
		h.Class(styles.ToggleSwitch),
		h.Input(onAttrs...),
		h.Input(offAttrs...),
		h.Label(
			h.For(field.Name+"_on"),
			h.Class(styles.Merge(styles.ToggleSwitchLabel, styles.ToggleSwitchLabelOn)),
			g.Text("On"),
		),
		h.Label(
			h.For(field.Name+"_off"),
			h.Class(styles.Merge(styles.ToggleSwitchLabel, styles.ToggleSwitchLabelOff)),
			g.Text("Off"),
		),
		h.Div(h.Class(styles.ToggleSwitchSlider)),
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

		attrs = appendFieldAttrs(attrs, field)

		if value != nil && currentVal == opt.Value {
			attrs = append(attrs, h.Checked())
		}

		options = append(options, h.Div(
			h.Class(styles.RadioOption),
			h.Input(attrs...),
			h.Label(
				h.For(field.Name+"_"+opt.Value),
				g.Text(opt.Label),
			),
		))
	}

	return h.Div(
		h.Class(styles.RadioGroup),
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
		h.Class(styles.FieldSelect),
	}

	attrs = appendFieldAttrs(attrs, field)

	return h.Select(
		g.Group(attrs),
		g.Group(options),
	)
}
