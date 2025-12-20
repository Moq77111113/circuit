package layout

import (
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
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
	result := Sidebar(s, values)
	html := renderToString(result)

	tests := []struct {
		name     string
		expected string
	}{
		{"root struct link", `href="#section-Server"`},
		{"root field link", `href="#field-Debug"`},
		{"nested struct link", `href="#section-Server.Database"`},
		{"nested field in root struct", `href="#field-Server.Host"`},
		{"nested field in nested struct", `href="#field-Server.Database.User"`},
		{"nested field in nested struct 2", `href="#field-Server.Database.Port"`},
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
	result := Sidebar(s, values)
	html := renderToString(result)

	if !strings.Contains(html, `href="#slice-Services"`) {
		t.Errorf("expected slice container link in HTML, got:\n%s", html)
	}

	if !strings.Contains(html, `href="#slice-item-Services.0"`) {
		t.Errorf("expected link to Services.0 in HTML, got:\n%s", html)
	}

	if !strings.Contains(html, `href="#slice-item-Services.1"`) {
		t.Errorf("expected link to Services.1 in HTML, got:\n%s", html)
	}
	if !strings.Contains(html, "Item 0") || !strings.Contains(html, "Item 1") {
		t.Errorf("expected Item 0 and Item 1 labels in HTML, got:\n%s", html)
	}
}

func TestSidebar_CollapsibleItems(t *testing.T) {
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
	result := Sidebar(s, values)
	html := renderToString(result)
	if !strings.Contains(html, "nav__item--collapsible") {
		t.Errorf("expected nav__item--collapsible class in HTML, got:\n%s", html)
	}

	if !strings.Contains(html, "nav__chevron") {
		t.Errorf("expected nav__chevron in HTML, got:\n%s", html)
	}

	if !strings.Contains(html, `onclick="toggleNavItem(this)"`) {
		t.Errorf("expected toggleNavItem onclick handler in HTML, got:\n%s", html)
	}
}
