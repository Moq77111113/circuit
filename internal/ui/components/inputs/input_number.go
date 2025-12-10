package inputs

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Number(field tags.Field, value any) g.Node {
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
	return h.Input(append(attrs, h.Type("number"))...)
}
