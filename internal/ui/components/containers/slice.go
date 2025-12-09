package containers

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Slice(field tags.Field, value any) g.Node {
	return SliceWithDepth(field, value, 0)
}

func SliceWithDepth(field tags.Field, value any, depth int) g.Node {
	items := reflectSliceValues(value)

	var nodes []g.Node

	if len(items) == 0 {
		nodes = append(nodes, h.Div(
			h.Class("slice__empty"),
			g.Text("No items"),
		))
	} else {
		for i, item := range items {
			nodes = append(nodes, RenderSliceItem(field, i, item, depth))
		}
	}

	nodes = append(nodes, h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("add:%s", field.Name)),
		h.Class("slice__add-btn"),
		g.Text("Add Item"),
	))

	header := CollapsibleHeader(field.Name, len(items), IsCollapsed(depth))
	body := CollapsibleBody(nodes)

	return CollapsibleContainer(depth, header, body)
}
