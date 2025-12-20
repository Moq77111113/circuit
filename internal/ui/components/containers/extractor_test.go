package containers

import (
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
)

func TestExtract_SingleString(t *testing.T) {
	node := ast.Node{
		Name: "Service",
		Kind: ast.KindStruct,
		Children: []ast.Node{
			{Name: "Name", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
		},
	}
	value := struct{ Name string }{Name: "User Service"}

	got := Extract(node, value, 3)

	if len(got.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(got.Fields))
	}
	if got.Fields[0].Name != "Name" {
		t.Errorf("expected Name, got %s", got.Fields[0].Name)
	}
	if got.Fields[0].Value != "User Service" {
		t.Errorf("expected 'User Service', got %s", got.Fields[0].Value)
	}
}

func TestExtract_MixedTypes(t *testing.T) {
	node := ast.Node{
		Name: "Service",
		Kind: ast.KindStruct,
		Children: []ast.Node{
			{Name: "Name", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
			{Name: "Enabled", Kind: ast.KindPrimitive, ValueType: ast.ValueBool},
			{Name: "Port", Kind: ast.KindPrimitive, ValueType: ast.ValueInt},
		},
	}
	value := struct {
		Name    string
		Enabled bool
		Port    int
	}{Name: "API", Enabled: true, Port: 8080}

	got := Extract(node, value, 3)

	if len(got.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(got.Fields))
	}
}

func TestExtract_NilValue(t *testing.T) {
	node := ast.Node{Name: "Test", Kind: ast.KindStruct}
	got := Extract(node, nil, 3)

	if len(got.Fields) != 0 {
		t.Errorf("expected 0 fields for nil value, got %d", len(got.Fields))
	}
}

func TestExtract_MaxFields(t *testing.T) {
	node := ast.Node{
		Name: "Service",
		Kind: ast.KindStruct,
		Children: []ast.Node{
			{Name: "Field1", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
			{Name: "Field2", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
			{Name: "Field3", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
			{Name: "Field4", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
		},
	}
	value := struct {
		Field1, Field2, Field3, Field4 string
	}{"A", "B", "C", "D"}

	got := Extract(node, value, 2)

	if len(got.Fields) != 2 {
		t.Errorf("expected max 2 fields, got %d", len(got.Fields))
	}
}

func TestExtract_ZeroValues(t *testing.T) {
	node := ast.Node{
		Name: "Service",
		Kind: ast.KindStruct,
		Children: []ast.Node{
			{Name: "Name", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
			{Name: "Enabled", Kind: ast.KindPrimitive, ValueType: ast.ValueBool},
			{Name: "Count", Kind: ast.KindPrimitive, ValueType: ast.ValueInt},
		},
	}
	value := struct {
		Name    string
		Enabled bool
		Count   int
	}{Name: "", Enabled: false, Count: 0}

	got := Extract(node, value, 3)

	if len(got.Fields) != 0 {
		t.Errorf("expected 0 fields (all zero values), got %d", len(got.Fields))
	}
}

func TestFormat(t *testing.T) {
	summary := Summary{
		Fields: []Field{
			{Name: "Name", Value: "User Service"},
			{Name: "Type", Value: "HTTP"},
		},
	}

	got := Format(summary)
	expected := "Name: User Service • Type: HTTP"

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormat_Empty(t *testing.T) {
	summary := Summary{Fields: []Field{}}
	got := Format(summary)

	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
