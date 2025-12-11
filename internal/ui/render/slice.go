package render

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/reflection"
	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/ui/components/containers"
)

// SliceRenderer renders slice fields
type SliceRenderer struct {
	dispatcher *Dispatcher
}

func (r *SliceRenderer) Render(node schema.Node, ctx Context) g.Node {
	items := reflection.SliceValues(ctx.Value)

	var itemNodes []g.Node
	if len(items) == 0 {
		itemNodes = append(itemNodes, r.renderEmptyState())
	} else {
		for i, itemValue := range items {
			itemNodes = append(itemNodes, r.renderItem(node, i, itemValue, ctx))
		}
	}

	itemNodes = append(itemNodes, r.renderAddButton(ctx.Path))

	isCollapsed := ctx.Depth >= 2
	header := containers.CollapsibleHeader(node.Name, len(items), isCollapsed, "")
	body := containers.CollapsibleBody(itemNodes)
	return containers.CollapsibleContainer(ctx.Depth, ctx.Path.String(), header, body)
}

func (r *SliceRenderer) renderEmptyState() g.Node {
	return h.Div(
		h.Class("slice__empty"),
		g.Text("No items"),
	)
}

func (r *SliceRenderer) renderAddButton(path schema.Path) g.Node {
	return h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("add:%s", path.String())),
		h.Class("slice__add-button"),
		g.Text("Add Item"),
	)
}

func (r *SliceRenderer) renderItem(node schema.Node, index int, value any, ctx Context) g.Node {
	itemPath := ctx.Path.Index(index)

	if node.ElementKind == schema.KindPrimitive {
		return r.renderPrimitiveItem(node, index, value, itemPath)
	}

	return r.renderStructItem(node, index, value, itemPath, ctx.Depth)
}
