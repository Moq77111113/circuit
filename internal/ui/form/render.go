package form

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/ui/render"
)

var dispatcher = render.NewDispatcher()

func renderNode(node schema.Node, values map[string]any) g.Node {
	var value any
	if values != nil {
		value = values[node.Name]
	}

	ctx := render.Context{
		Path:  schema.NewPath(node.Name),
		Value: value,
		Depth: 0,
	}

	return dispatcher.Render(node, ctx)
}
