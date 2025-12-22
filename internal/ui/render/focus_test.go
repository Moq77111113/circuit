package render

import (
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

func TestFilterByFocus_Root(t *testing.T) {
	nodes := []ast.Node{
		{Name: "AppName", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
		{Name: "Port", Kind: ast.KindPrimitive, ValueType: ast.ValueInt},
		{
			Name: "Database",
			Kind: ast.KindStruct,
			Children: []ast.Node{
				{Name: "Host", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
				{Name: "Port", Kind: ast.KindPrimitive, ValueType: ast.ValueInt},
			},
		},
	}

	result := FilterByFocus(nodes, path.Root())

	// Root should show primitives + struct as card (no children expanded)
	if len(result) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(result))
	}

	// Check struct node has no children at root level (will be rendered as card)
	dbNode := result[2]
	if dbNode.Name != "Database" {
		t.Errorf("expected Database node, got %s", dbNode.Name)
	}
}

func TestFilterByFocus_Struct(t *testing.T) {
	nodes := []ast.Node{
		{Name: "AppName", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
		{
			Name: "Database",
			Kind: ast.KindStruct,
			Children: []ast.Node{
				{Name: "Host", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
				{Name: "Port", Kind: ast.KindPrimitive, ValueType: ast.ValueInt},
			},
		},
	}

	focusPath := path.Root().Child("Database")
	result := FilterByFocus(nodes, focusPath)

	// Should return children of Database
	if len(result) != 2 {
		t.Fatalf("expected 2 nodes (Database children), got %d", len(result))
	}

	if result[0].Name != "Host" {
		t.Errorf("expected Host, got %s", result[0].Name)
	}
	if result[1].Name != "Port" {
		t.Errorf("expected Port, got %s", result[1].Name)
	}
}

func TestFilterByFocus_SliceItem(t *testing.T) {
	nodes := []ast.Node{
		{
			Name:        "Services",
			Kind:        ast.KindSlice,
			ElementKind: ast.KindStruct,
			Children: []ast.Node{
				{Name: "Name", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
				{Name: "Port", Kind: ast.KindPrimitive, ValueType: ast.ValueInt},
			},
		},
	}

	focusPath := path.Root().Child("Services").Index(0)
	result := FilterByFocus(nodes, focusPath)

	// Should return children of Services[0]
	if len(result) != 2 {
		t.Fatalf("expected 2 nodes (slice item children), got %d", len(result))
	}

	if result[0].Name != "Name" {
		t.Errorf("expected Name, got %s", result[0].Name)
	}
}

func TestFilterByFocus_DeepNesting(t *testing.T) {
	nodes := []ast.Node{
		{
			Name: "Server",
			Kind: ast.KindStruct,
			Children: []ast.Node{
				{
					Name: "RateLimit",
					Kind: ast.KindStruct,
					Children: []ast.Node{
						{Name: "Enabled", Kind: ast.KindPrimitive, ValueType: ast.ValueBool},
						{Name: "Rate", Kind: ast.KindPrimitive, ValueType: ast.ValueInt},
					},
				},
			},
		},
	}

	focusPath := path.Root().Child("Server").Child("RateLimit")
	result := FilterByFocus(nodes, focusPath)

	// Should return children of Server.RateLimit
	if len(result) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(result))
	}

	if result[0].Name != "Enabled" {
		t.Errorf("expected Enabled, got %s", result[0].Name)
	}
}
