package inputs

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

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
