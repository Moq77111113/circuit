package form

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ui/render"
	"github.com/moq77111113/circuit/internal/ui/styles"
)

func Form(s ast.Schema, values path.ValuesByPath, focus path.Path) g.Node {
	filteredNodes := render.FilterByFocus(s.Nodes, focus)
	basePath := computeBasePath(filteredNodes, focus)
	fields := render.Render(filteredNodes, values, basePath)

	return h.Form(
		h.Method("post"),
		h.Class(styles.Form),
		fields,
		h.Div(
			h.Class(styles.FormActions),
			h.Button(
				h.Type("submit"),
				h.Class(styles.Button+" "+styles.ButtonPrimary),
				g.Text("Save Changes"),
			),
		),
	)
}

func computeBasePath(nodes []ast.Node, focus path.Path) path.Path {
	if focus.IsRoot() {
		return path.Root()
	}
	if len(nodes) == 1 {
		segments := focus.Segments()
		if len(segments) > 0 && nodes[0].Name == segments[len(segments)-1] {

			parentPath := path.Root()
			for _, seg := range segments[:len(segments)-1] {
				parentPath = parentPath.Child(seg)
			}
			return parentPath
		}
	}
	return focus
}

func renderToString(node g.Node) string {
	var sb strings.Builder
	_ = node.Render(&sb)
	return sb.String()
}
