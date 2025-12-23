package render

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

// parseItemPath extracts the field path and index from a full item path
func parseItemPath(itemPath string) (field string, index string) {
	lastDot := strings.LastIndex(itemPath, ".")
	if lastDot == -1 {
		return "", itemPath
	}
	return itemPath[:lastDot], itemPath[lastDot+1:]
}

// renderAddButton creates an "Add" button for slices
func renderAddButton(path path.Path) g.Node {
	return h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("add:%s", path.FieldPath())),
		h.Class("button button--primary"),
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
	field, idx := parseItemPath(itemPath)
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
			h.Value(fmt.Sprintf("remove:%s:%s", field, idx)),
			h.Class("button button--secondary"+" slice-item__remove-button"),
			g.Text("Remove"),
		),
	)
}
