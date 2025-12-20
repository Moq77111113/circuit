package walk

import (
	"reflect"
	"testing"

	"github.com/moq77111113/circuit/internal/ast/node"
)

// MockVisitor for testing
type MockVisitor struct {
	visited []string
}

func (m *MockVisitor) VisitPrimitive(ctx *VisitContext, node *node.Node) error {
	m.visited = append(m.visited, node.Name)
	return nil
}

func (m *MockVisitor) VisitStruct(ctx *VisitContext, node *node.Node) error {
	m.visited = append(m.visited, node.Name)
	return nil
}

func (m *MockVisitor) VisitSlice(ctx *VisitContext, node *node.Node) error {
	m.visited = append(m.visited, node.Name)
	return nil
}

// PathRecorderVisitor records paths during traversal
type PathRecorderVisitor struct {
	paths []string
}

func (p *PathRecorderVisitor) VisitPrimitive(ctx *VisitContext, node *node.Node) error {
	p.paths = append(p.paths, ctx.Path.String())
	return nil
}

func (p *PathRecorderVisitor) VisitStruct(ctx *VisitContext, node *node.Node) error {
	p.paths = append(p.paths, ctx.Path.String())
	return nil
}

func (p *PathRecorderVisitor) VisitSlice(ctx *VisitContext, node *node.Node) error {
	p.paths = append(p.paths, ctx.Path.String())
	return nil
}

func TestWalker_VisitsAllNodes(t *testing.T) {
	tree := &node.Tree{Nodes: []node.Node{
		{Name: "Field1", Kind: node.KindPrimitive},
		{Name: "Field2", Kind: node.KindStruct, Children: []node.Node{
			{Name: "Nested", Kind: node.KindPrimitive},
		}},
	}}

	visitor := &MockVisitor{}
	walker := NewWalker(visitor)
	err := walker.Walk(tree, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"Field1", "Field2", "Nested"}
	if !reflect.DeepEqual(visitor.visited, expected) {
		t.Errorf("got %v, want %v", visitor.visited, expected)
	}
}

func TestWalker_RespectsMaxDepth(t *testing.T) {
	tree := &node.Tree{Nodes: []node.Node{
		{Name: "Level0", Kind: node.KindStruct, Children: []node.Node{
			{Name: "Level1", Kind: node.KindStruct, Children: []node.Node{
				{Name: "Level2", Kind: node.KindPrimitive},
			}},
		}},
	}}

	visitor := &MockVisitor{}
	walker := NewWalker(visitor, WithMaxDepth(1))
	walker.Walk(tree, nil)

	// Level2 should NOT be visited (depth = 2)
	for _, name := range visitor.visited {
		if name == "Level2" {
			t.Error("Level2 should not be visited with MaxDepth=1")
		}
	}
}

func TestWalker_TracksPath(t *testing.T) {
	tree := &node.Tree{Nodes: []node.Node{
		{Name: "Server", Kind: node.KindStruct, Children: []node.Node{
			{Name: "Port", Kind: node.KindPrimitive},
		}},
	}}

	visitor := &PathRecorderVisitor{}
	walker := NewWalker(visitor)
	walker.Walk(tree, nil)

	expected := []string{"Server", "Server.Port"}
	if !reflect.DeepEqual(visitor.paths, expected) {
		t.Errorf("got %v, want %v", visitor.paths, expected)
	}
}

func TestWalker_EmptyTree(t *testing.T) {
	tree := &node.Tree{Nodes: []node.Node{}}
	visitor := &MockVisitor{}
	walker := NewWalker(visitor)

	err := walker.Walk(tree, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(visitor.visited) != 0 {
		t.Errorf("expected no visits, got %d", len(visitor.visited))
	}
}
