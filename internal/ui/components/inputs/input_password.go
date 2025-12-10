package inputs

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Password(field tags.Field, value any) g.Node {
	attrs := BaseAttrs(field)
	if value != nil {
		attrs = append(attrs, h.Value(fmt.Sprintf("%v", value)))
	}
	return h.Input(append(attrs, h.Type("password"))...)
}
