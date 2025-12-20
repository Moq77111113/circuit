package node

import "testing"

func TestTree_Find(t *testing.T) {
	tree := &Tree{Nodes: []Node{
		{ID: "Config.Server", Name: "Server", Children: []Node{
			{ID: "Config.Server.Port", Name: "Port"},
		}},
	}}

	node := tree.Find("Config.Server.Port")
	if node == nil {
		t.Fatal("expected to find node")
	}
	if node.Name != "Port" {
		t.Errorf("got %s, want Port", node.Name)
	}
}

func TestTree_Find_NotFound(t *testing.T) {
	tree := &Tree{Nodes: []Node{}}
	node := tree.Find("DoesNotExist")
	if node != nil {
		t.Error("expected nil")
	}
}

func TestTree_Find_Nested(t *testing.T) {
	tree := &Tree{Nodes: []Node{
		{
			ID:   "Config.Server",
			Name: "Server",
			Children: []Node{
				{
					ID:   "Config.Server.Database",
					Name: "Database",
					Children: []Node{
						{ID: "Config.Server.Database.Host", Name: "Host"},
					},
				},
			},
		},
	}}

	node := tree.Find("Config.Server.Database.Host")
	if node == nil {
		t.Fatal("expected to find deeply nested node")
	}
	if node.Name != "Host" {
		t.Errorf("got %s, want Host", node.Name)
	}
}
