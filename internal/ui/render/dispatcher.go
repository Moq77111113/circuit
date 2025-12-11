package render

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/schema"
)

type Context struct {
	Path  schema.Path // Current field path (e.g., "Services.0.Endpoints.2.Name")
	Value any         // Current field value from the config
	Depth int         // Nesting depth (used for collapsed rendering at depth >= 2)
}


type Renderer interface {
	Render(node schema.Node, ctx Context) g.Node
}

type Dispatcher struct {
	primitive  Renderer
	structNode Renderer
	sliceNode  Renderer
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{}
	d.primitive = &PrimitiveRenderer{}
	d.structNode = &StructRenderer{dispatcher: d}
	d.sliceNode = &SliceRenderer{dispatcher: d}
	return d
}

func (d *Dispatcher) Render(node schema.Node, ctx Context) g.Node {
	switch node.Kind {
	case schema.KindPrimitive:
		return d.primitive.Render(node, ctx)
	case schema.KindStruct:
		return d.structNode.Render(node, ctx)
	case schema.KindSlice:
		return d.sliceNode.Render(node, ctx)
	default:
		return nil
	}
}
