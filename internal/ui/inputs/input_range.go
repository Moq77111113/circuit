package inputs

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Range(field tags.Field, value any) g.Node {
	attrs := BaseAttrs(field)
	if value != nil {
		attrs = append(attrs, h.Value(fmt.Sprintf("%v", value)))
	}
	if field.Min != "" {
		attrs = append(attrs, h.Min(field.Min))
	}
	if field.Max != "" {
		attrs = append(attrs, h.Max(field.Max))
	}
	if field.Step != "" {
		attrs = append(attrs, h.Step(field.Step))
	}

	attrs = append(attrs, g.Attr("oninput", "this.nextElementSibling.nextElementSibling.value = this.value"))

	return h.Div(
		h.Class("range-wrapper"),
		h.Span(h.Class("range-min"), g.Text(field.Min)),
		h.Input(append(attrs, h.Type("range"))...),
		h.Span(h.Class("range-max"), g.Text(field.Max)),
		g.El("output", h.Class("range-value"), g.Text(fmt.Sprintf("%v", value))),
	)
}
