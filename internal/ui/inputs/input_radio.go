package inputs

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

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
