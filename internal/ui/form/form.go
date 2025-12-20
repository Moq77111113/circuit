package form

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ui/render"
)

func Form(s ast.Schema, values map[string]any) g.Node {
	// Use the new Render API that takes all nodes at once
	fields := render.Render(s.Nodes, values)

	return h.Form(
		h.Method("post"),
		h.Class("form"),
		fields,
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
