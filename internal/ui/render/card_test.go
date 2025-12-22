package render

import (
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

func TestRenderStructCard(t *testing.T) {
	node := ast.Node{
		Name: "Database",
		Kind: ast.KindStruct,
		Children: []ast.Node{
			{Name: "Host", Kind: ast.KindPrimitive, ValueType: ast.ValueString},
			{Name: "Port", Kind: ast.KindPrimitive, ValueType: ast.ValueInt},
		},
	}

	values := map[string]any{
		"Database.Host": "localhost",
		"Database.Port": 5432,
	}

	nodePath := path.Root().Child("Database")
	html := RenderStructCard(node, nodePath, values)

	result := renderToString(html)

	if !strings.Contains(result, "Database") {
		t.Error("card should contain struct name")
	}

	if !strings.Contains(result, "localhost") {
		t.Error("card should contain host value in preview")
	}

	if !strings.Contains(result, "5432") {
		t.Error("card should contain port value in preview")
	}

	if !strings.Contains(result, "?focus=Database") {
		t.Error("card should link to focus URL")
	}
}

func TestRenderStructCard_DeepNesting(t *testing.T) {
	node := ast.Node{
		Name: "RateLimit",
		Kind: ast.KindStruct,
		Children: []ast.Node{
			{Name: "Enabled", Kind: ast.KindPrimitive, ValueType: ast.ValueBool},
			{Name: "Rate", Kind: ast.KindPrimitive, ValueType: ast.ValueInt},
		},
	}

	values := map[string]any{
		"Server.RateLimit.Enabled": true,
		"Server.RateLimit.Rate":    100,
	}

	nodePath := path.Root().Child("Server").Child("RateLimit")
	html := RenderStructCard(node, nodePath, values)

	result := renderToString(html)

	if !strings.Contains(result, "?focus=Server.RateLimit") {
		t.Errorf("card should link to deep focus path, got: %s", result)
	}
}
