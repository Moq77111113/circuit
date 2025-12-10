package containers

import (
	"fmt"
	"strings"

	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Render(field tags.Field, index int, value any, depth int) g.Node {
	summary := Extract(field, value, 3)
	summaryText := formatCompact(summary)

	removebutton := h.Button(
		h.Type("submit"),
		h.Name("action"),
		h.Value(fmt.Sprintf("remove:%s:%d", field.Name, index)),
		h.Class("slice__remove-button"),
		g.Text("Remove"),
	)

	return h.Div(
		h.Class("slice__item slice__item--compact collapsed"),
		h.Div(
			h.Class("slice__item-header"),
			g.Attr("onclick", "toggleCollapse(this)"),
			g.Text(fmt.Sprintf("#%d", index+1)),
			g.Text(" | "),
			g.Text(summaryText),
			g.Text(" ▼ "),
			removebutton,
		),
		h.Div(
			h.Class("slice__item-body"),
			renderStructItem(field, index, value, depth),
		),
	)
}

func formatCompact(s Summary) string {
	if len(s.Fields) == 0 {
		return ""
	}

	var parts []string
	for _, f := range s.Fields {
		parts = append(parts, fmt.Sprintf("%s: %s", f.Name, f.Value))
	}
	return strings.Join(parts, " | ")
}
