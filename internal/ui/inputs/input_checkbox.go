package inputs

import (
	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Checkbox(field tags.Field, value any) g.Node {
	attrs := BaseAttrs(field)
	if value != nil && value.(bool) {
		attrs = append(attrs, h.Checked())
	}
	return h.Input(append(attrs, h.Type("checkbox"))...)
}
