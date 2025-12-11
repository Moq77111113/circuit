package render

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/reflection"
	"github.com/moq77111113/circuit/internal/schema"
)

// renderStructItem renders a struct slice item
func (r *SliceRenderer) renderStructItem(node schema.Node, index int, value any, path schema.Path, depth int) g.Node {
	removeBtn := h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("remove:%s:%d", path.String(), index)),
		h.Class("slice__remove-button"),
		g.Text("Remove"),
	)

	var fieldNodes []g.Node
	for _, child := range node.Children {
		childCtx := Context{
			Path:  path.Child(child.Name),
			Value: reflection.FieldByName(value, child.Name),
			Depth: depth + 1,
		}
		fieldNodes = append(fieldNodes, r.dispatcher.Render(child, childCtx))
	}

	if depth >= 2 {
		return r.renderCompactStructItem(node, index, value, path, fieldNodes, removeBtn)
	}

	return r.renderStandardStructItem(node, index, value, path, fieldNodes, removeBtn)
}

func (r *SliceRenderer) renderCompactStructItem(node schema.Node, index int, value any, path schema.Path, fieldNodes []g.Node, removeButton g.Node) g.Node {
	summary := extractSummary(node, value, 3)
	summaryText := formatSummary(summary)

	return h.Div(
		h.ID("slice-item-"+path.String()),
		h.Class("slice__item slice__item--compact collapsed"),
		h.Div(
			h.Class("slice__item-header"),
			g.Attr("onclick", "toggleCollapse(this)"),
			g.Text(fmt.Sprintf("#%d | %s ▾ ", index+1, summaryText)),
			removeButton,
		),
		h.Div(
			h.Class("slice__item-body"),
			g.Group(fieldNodes),
		),
	)
}

func (r *SliceRenderer) renderStandardStructItem(node schema.Node, index int, value any, path schema.Path, fieldNodes []g.Node, removeBtn g.Node) g.Node {
	summary := extractSummary(node, value, 3)
	title := fmt.Sprintf("#%d", index+1)
	if len(summary.Fields) > 0 {
		title = fmt.Sprintf("%s: %s", summary.Fields[0].Name, summary.Fields[0].Value)
	}

	return h.Div(
		h.ID("slice-item-"+path.String()),
		h.Class("slice__item slice__item--struct"),
		h.Div(
			h.Class("slice__item-header"),
			g.Attr("onclick", "toggleCollapse(this)"),
			h.Span(h.Class("slice__chevron"), g.Text("▾")),
			h.Span(h.Class("slice__item-title"), g.Text(title)),
			removeBtn,
		),
		h.Div(
			h.Class("slice__item-body"),
			g.Group(fieldNodes),
		),
	)
}
