package inputs

import (
	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func BaseAttrs(field tags.Field) []g.Node {
	attrs := []g.Node{
		h.Name(field.Name),
		h.ID(field.Name),
		h.Class("field__input"),
	}

	if field.Required {
		attrs = append(attrs, h.Required())
	}
	return attrs
}
