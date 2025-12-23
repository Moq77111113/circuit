package layout

import (
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

func TestSidebar_RootLevelOnly(t *testing.T) {
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

	rootTests := []struct {
		name     string
		expected string
	}{
		{"tree structure", `class="sidebar-tree"`},
		{"root label", `Config`},
		{"root struct link", `href="?focus=Server"`},
		{"root field link", `href="?focus=Debug"`},
		{"tree node class", `class="tree-node`},
	}

	for _, tt := range rootTests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(html, tt.expected) {
				t.Errorf("expected %q in HTML, got:\n%s", tt.expected, html)
			}
		})
	}

	nestedTests := []struct {
		name       string
		unexpected string
	}{
		{"no nested struct", `href="?focus=Server.Database"`},
		{"no nested field in root", `href="?focus=Server.Host"`},
		{"no nested field in nested struct", `href="?focus=Server.Database.User"`},
	}

	for _, tt := range nestedTests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.Contains(html, tt.unexpected) {
				t.Errorf("should NOT contain %q in HTML (only root-level shown), got:\n%s", tt.unexpected, html)
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

	// Only root-level elements should be present
	tests := []struct {
		name     string
		expected string
	}{
		{"tree structure", `class="sidebar-tree"`},
		{"tree node", `class="tree-node`},
		{"server link", `href="?focus=Server"`},
		{"tags link", `href="?focus=Tags"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(html, tt.expected) {
				t.Errorf("expected %q in HTML, got:\n%s", tt.expected, html)
			}
		})
	}

	// Nested elements should NOT be present
	if strings.Contains(html, `href="?focus=Server.Host"`) {
		t.Errorf("should NOT contain nested Host link (only root-level shown)")
	}
}
