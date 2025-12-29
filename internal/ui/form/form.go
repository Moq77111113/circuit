package form

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ui/render"
	"github.com/moq77111113/circuit/internal/ui/styles"
	"github.com/moq77111113/circuit/internal/validation"
)

func Form(s ast.Schema, values ast.ValuesByPath, focus path.Path, readOnly bool) g.Node {
	filteredNodes := render.FilterByFocus(s.Nodes, focus)
	basePath := computeBasePath(filteredNodes, focus)

	opts := render.DefaultOptions()
	opts.ReadOnly = readOnly
	fields := render.RenderWithOptions(filteredNodes, values, basePath, opts)

	var actions g.Node
	if !readOnly {
		actions = h.Div(
			h.Class(styles.FormActions),
			h.Button(
				h.Type("submit"),
				h.Class(styles.Button+" "+styles.ButtonPrimary),
				g.Text("Save Changes"),
			),
		)
	}

	return h.Form(
		h.Method("post"),
		h.Class(styles.Form),
		fields,
		actions,
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

// FormWithErrors renders the form with validation errors.
func FormWithErrors(s ast.Schema, values ast.ValuesByPath, focus path.Path, errors *validation.ValidationResult, readOnly bool) g.Node {
	filteredNodes := render.FilterByFocus(s.Nodes, focus)
	basePath := computeBasePath(filteredNodes, focus)

	opts := render.DefaultOptions()
	opts.Errors = errors
	opts.ReadOnly = readOnly

	fields := render.RenderWithOptions(filteredNodes, values, basePath, opts)

	var actions g.Node
	if !readOnly {
		actions = h.Div(
			h.Class(styles.FormActions),
			h.Button(
				h.Type("submit"),
				h.Class(styles.Button+" "+styles.ButtonPrimary),
				g.Text("Save Changes"),
			),
		)
	}

	return h.Form(
		h.Method("post"),
		h.Class(styles.Form),
		fields,
		actions,
	)
}
