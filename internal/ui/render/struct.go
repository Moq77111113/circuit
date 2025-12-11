package render

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/reflection"
	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/ui/components/containers"
)

type StructRenderer struct {
	dispatcher *Dispatcher
}

// Render generates a section container with all child fields.
func (r *StructRenderer) Render(node schema.Node, ctx Context) g.Node {
	var children []g.Node
	for _, child := range node.Children {
		childCtx := Context{
			Path:  ctx.Path.Child(child.Name),
			Value: reflection.FieldByName(ctx.Value, child.Name),
			Depth: ctx.Depth + 1,
		}
		children = append(children, r.dispatcher.Render(child, childCtx))
	}

	return containers.Section(ctx.Path.String(), children, true, false)
}
