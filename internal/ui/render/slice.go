package render

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

// renderAddButton creates an "Add" button for slices
func renderAddButton(path path.Path) g.Node {
	return h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("add_%s", path.String())),
		h.Class("btn btn--add"),
		g.Text("Add"),
	)
}

// renderEmptyState returns a message for empty slices
func renderEmptyState() g.Node {
	return h.P(
		h.Class("empty-state"),
		g.Text("No items"),
	)
}

// renderPrimitiveSliceItem renders a single primitive item in a slice
func renderPrimitiveSliceItem(node *ast.Node, index int, value any, path path.Path) g.Node {
	itemPath := path.String()
	return h.Div(
		h.Class("slice-item slice-item--primitive"),
		h.Div(
			h.Class("field"),
			renderLabel(node, itemPath),
			renderInput(node, itemPath, value),
		),
		h.Button(
			h.Type("submit"),
			h.Name("action"),
			h.Value(fmt.Sprintf("remove_%s", itemPath)),
			h.Class("btn btn--remove"),
			g.Text("Remove"),
		),
	)
}

// renderStructSliceItem renders a single struct item in a slice (collapsed summary)
func renderStructSliceItem(index int, summary string, path path.Path) g.Node {
	return h.Div(
		h.Class("slice-item slice-item--struct"),
		h.A(
			h.Href("#"+path.String()),
			h.Class("slice-item__link"),
			g.Textf("#%d: %s", index, summary),
		),
		h.Button(
			h.Type("submit"),
			h.Name("action"),
			h.Value(fmt.Sprintf("remove_%s", path.String())),
			h.Class("btn btn--remove"),
			g.Text("Remove"),
		),
	)
}
