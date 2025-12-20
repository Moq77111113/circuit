package sidebar

import (
	g "maragu.dev/gomponents"
)

// LinkState accumulates navigation links during tree traversal.
type LinkState struct {
	links []g.Node
}

// NewLinkState creates a new link accumulation state.
func NewLinkState() *LinkState {
	return &LinkState{
		links: make([]g.Node, 0),
	}
}

// Add appends a link to the accumulated links.
func (s *LinkState) Add(link g.Node) {
	s.links = append(s.links, link)
}

// Links returns all accumulated links.
func (s *LinkState) Links() []g.Node {
	return s.links
}
