package node

// Tree represents a complete AST with bidirectional navigation.
type Tree struct {
	Nodes []Node
}

// Find returns the node with the given ID, or nil if not found.
func (t *Tree) Find(id string) *Node {
	return findInNodes(t.Nodes, id)
}

func findInNodes(nodes []Node, id string) *Node {
	for i := range nodes {
		if nodes[i].ID == id {
			return &nodes[i]
		}
		if len(nodes[i].Children) > 0 {
			if found := findInNodes(nodes[i].Children, id); found != nil {
				return found
			}
		}
	}
	return nil
}
