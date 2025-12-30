package render

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ui/styles"
)

// parseItemPath extracts the field path and index from a full item path
func parseItemPath(itemPath string) (field string, index string) {
	lastDot := strings.LastIndex(itemPath, ".")
	if lastDot == -1 {
		return "", itemPath
	}
	return itemPath[:lastDot], itemPath[lastDot+1:]
}

// renderAddButton creates an "Add" button for slices (returns nil if readOnly)
func renderAddButton(path path.Path, readOnly bool) g.Node {
	if readOnly {
		return nil
	}
	return h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("add:%s", path.FieldPath())),
		h.Class(styles.Merge(styles.Button, styles.ButtonPrimary, styles.ButtonAdd)),
		g.Text("Add"),
	)
}

// renderEmptyState returns a message for empty slices
func renderEmptyState() g.Node {
	return h.P(
		h.Class(styles.EmptyState),
		g.Text("No items"),
	)
}

// renderPrimitiveSliceItem renders a single primitive item in a slice
func renderPrimitiveSliceItem(node *ast.Node, index int, value any, path path.Path, rc *RenderContext) g.Node {
	itemPath := path.String()
	field, idx := parseItemPath(itemPath)

	var removeBtn g.Node
	if !rc.ReadOnly {
		removeBtn = h.Button(
			h.Type("submit"),
			h.Name("action"),
			h.Value(fmt.Sprintf("remove:%s:%s", field, idx)),
			h.Class(styles.Merge(styles.Button, styles.ButtonSecondary, styles.ButtonRemove, styles.SliceItemRemoveButton)),
			g.Text("Remove"),
		)
	}

	return h.Div(
		h.Class(styles.Merge(styles.SliceItem, styles.SliceItemPrimitive)),
		h.Div(
			h.Class(styles.Field),
			renderLabel(node, itemPath),
			renderInput(node, itemPath, value, rc),
		),
		removeBtn,
	)
}
