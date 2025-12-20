package render

import (
	g "maragu.dev/gomponents"
)

// RenderState accumulates rendered HTML nodes during tree traversal.
type RenderState struct {
	nodes []g.Node
}

// NewRenderState creates a new render state.
func NewRenderState() *RenderState {
	return &RenderState{
		nodes: make([]g.Node, 0),
	}
}

// Append adds a rendered node to the output.
func (s *RenderState) Append(node g.Node) {
	s.nodes = append(s.nodes, node)
}

// Output returns all accumulated nodes.
func (s *RenderState) Output() []g.Node {
	return s.nodes
}
