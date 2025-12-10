package containers

import (
	"fmt"

	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/components/inputs"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func RenderSliceItem(field tags.Field, index int, value any, depth int) g.Node {
	isStruct := len(field.Fields) > 0
	removebutton := h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("remove:%s:%d", field.Name, index)),
		h.Class("slice__remove-button"),
		g.Text("Remove"),
	)

	if isStruct {
		summary := Extract(field, value, 3)
		title := fmt.Sprintf("#%d", index+1)
		if len(summary.Fields) > 0 {
			title = fmt.Sprintf("%s: %s", summary.Fields[0].Name, summary.Fields[0].Value)
		}

		var summaryText string
		if len(summary.Fields) > 1 {
			summaryText = Format(Summary{Fields: summary.Fields[1:]})
		}

		titleNode := h.Span(h.Class("slice__item-title"), g.Text(title))
		var headerChildren []g.Node
		headerChildren = append(headerChildren,
			h.Span(h.Class("slice__chevron"), g.Text("▼")),
			titleNode,
		)
		if summaryText != "" {
			headerChildren = append(headerChildren,
				h.Span(h.Class("slice__summary"), g.Text(summaryText)),
			)
		}
		headerChildren = append(headerChildren, removebutton)

		return h.Div(
			h.Class("slice__item slice__item--struct"),
			h.Div(
				h.Class("slice__item-header"),
				g.Attr("onclick", "toggleCollapse(this)"),
				g.Group(headerChildren),
			),
			h.Div(
				h.Class("slice__item-body"),
				renderStructItem(field, index, value, depth),
			),
		)
	}

	return h.Div(
		h.Class("slice__item slice__item--primitive"),
		renderPrimitiveItem(field, index, value),
		removebutton,
	)
}

func renderPrimitiveItem(field tags.Field, index int, value any) g.Node {
	indexedName := fmt.Sprintf("%s.%d", field.Name, index)
	indexedField := tags.Field{
		Name:      indexedName,
		Type:      field.ElementType,
		InputType: field.InputType,
		Required:  field.Required,
		Min:       field.Min,
		Max:       field.Max,
		Step:      field.Step,
	}

	switch field.InputType {
	case tags.TypeText:
		return inputs.Text(indexedField, value)
	case tags.TypeNumber:
		return inputs.Number(indexedField, value)
	case tags.TypeCheckbox:
		return inputs.Checkbox(indexedField, value)
	default:
		return inputs.Text(indexedField, value)
	}
}
