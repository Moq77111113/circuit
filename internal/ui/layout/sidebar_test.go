package layout

import (
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

func TestSidebar_NestedStructLinks(t *testing.T) {
	nodes := []ast.Node{
		{
			Name: "Server",
			Kind: ast.KindStruct,
			Children: []ast.Node{
				{Name: "Host", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
				{
					Name: "Database",
					Kind: ast.KindStruct,
					Children: []ast.Node{
						{Name: "User", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
						{Name: "Port", Kind: ast.KindPrimitive, ValueType: ast.ValueInt},
					},
				},
			},
		},
		{Name: "Debug", Kind: ast.KindPrimitive, ValueType: ast.ValueBool},
	}

	s := ast.Schema{Nodes: nodes}
	values := map[string]any{
		"Server.Host":          "localhost",
		"Server.Database.User": "admin",
		"Server.Database.Port": 5432,
		"Debug":                true,
	}
	result := Sidebar(s, values, path.Root())
	html := renderToString(result)

	tests := []struct {
		name     string
		expected string
	}{
		{"tree structure", `class="sidebar-tree"`},
		{"root label", `Config`},
		{"root struct link", `href="?focus=Server"`},
		{"root field link", `href="?focus=Debug"`},
		{"nested struct link", `href="?focus=Server.Database"`},
		{"nested field in root struct", `href="?focus=Server.Host"`},
		{"nested field in nested struct", `href="?focus=Server.Database.User"`},
		{"nested field in nested struct 2", `href="?focus=Server.Database.Port"`},
		{"tree node class", `class="tree-node`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(html, tt.expected) {
				t.Errorf("expected %q in HTML, got:\n%s", tt.expected, html)
			}
		})
	}
}

func TestSidebar_SliceLinks(t *testing.T) {
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

	s := ast.Schema{Nodes: nodes}

	values := map[string]any{
		"Services": []map[string]any{
			{"Name": "API", "Port": 8080},
			{"Name": "DB", "Port": 5432},
		},
	}
	result := Sidebar(s, values, path.Root())
	html := renderToString(result)

	if !strings.Contains(html, `href="?focus=Services"`) {
		t.Errorf("expected slice link in HTML, got:\n%s", html)
	}

	if !strings.Contains(html, "Services") {
		t.Errorf("expected Services label in HTML, got:\n%s", html)
	}
}

func TestSidebar_TreeNavigation(t *testing.T) {
	nodes := []ast.Node{
		{
			Name: "Server",
			Kind: ast.KindStruct,
			Children: []ast.Node{
				{Name: "Host", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
			},
		},
		{
			Name:        "Tags",
			Kind:        ast.KindSlice,
			ElementKind: ast.KindPrimitive,
			ValueType:   ast.ValueString,
		},
	}

	s := ast.Schema{Nodes: nodes}
	values := map[string]any{
		"Server.Host": "localhost",
		"Tags":        []string{"prod", "api"},
	}
	result := Sidebar(s, values, path.Root())
	html := renderToString(result)

	tests := []struct {
		name     string
		expected string
	}{
		{"tree structure", `class="sidebar-tree"`},
		{"tree node", `class="tree-node`},
		{"server link", `href="?focus=Server"`},
		{"tags link", `href="?focus=Tags"`},
		{"host link", `href="?focus=Server.Host"`},
		{"chevron present", `▼`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(html, tt.expected) {
				t.Errorf("expected %q in HTML, got:\n%s", tt.expected, html)
			}
		})
	}
}
