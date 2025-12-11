package sidebar

import (
	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/schema"
)

// renderLinks generates navigation links for a list of nodes.
func RenderLinks(nodes []schema.Node, parentPath schema.Path, values map[string]any) []g.Node {
	var links []g.Node
	for _, node := range nodes {
		path := parentPath.Child(node.Name)

		switch node.Kind {
		case schema.KindStruct:
			links = append(links, renderStructLink(node, path, values))
		case schema.KindSlice:
			links = append(links, renderSliceLink(node, path, values))
		case schema.KindPrimitive:
			links = append(links, renderPrimitiveLink(node, path))
		}
	}
	return links
}
