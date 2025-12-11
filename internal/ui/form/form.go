package form

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/schema"
)

func Form(s schema.Schema, values map[string]any) g.Node {
	var fields []g.Node

	for _, node := range s.Nodes {
		fields = append(fields, renderNode(node, values))
	}

	return h.Form(
		h.Method("post"),
		h.Class("form"),
		g.Group(fields),
		h.Div(
			h.Class("form__actions"),
			h.Button(
				h.Type("submit"),
				h.Class("button button--primary"),
				g.Text("Save Changes"),
			),
		),
	)
}

func renderToString(node g.Node) string {
	var sb strings.Builder
	_ = node.Render(&sb)
	return sb.String()
}
