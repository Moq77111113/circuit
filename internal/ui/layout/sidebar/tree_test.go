package sidebar

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

func renderToString(node g.Node) string {
	var sb strings.Builder
	_ = node.Render(&sb)
	return sb.String()
}

func TestRenderTree_Simple(t *testing.T) {
	nodes := []ast.Node{
		{Name: "AppName", Kind: ast.KindPrimitive},
		{Name: "Port", Kind: ast.KindPrimitive},
	}

	html := RenderTree(nodes, path.Root(), nil)
	result := renderToString(html)

	if !strings.Contains(result, "AppName") {
		t.Error("tree should contain AppName")
	}

	if !strings.Contains(result, "Port") {
		t.Error("tree should contain Port")
	}
}

func TestRenderTree_WithStruct(t *testing.T) {
	nodes := []ast.Node{
		{Name: "AppName", Kind: ast.KindPrimitive},
		{
			Name: "Database",
			Kind: ast.KindStruct,
			Children: []ast.Node{
				{Name: "Host", Kind: ast.KindPrimitive},
				{Name: "Port", Kind: ast.KindPrimitive},
			},
		},
	}

	html := RenderTree(nodes, path.Root(), nil)
	result := renderToString(html)

	if !strings.Contains(result, "Database") {
		t.Error("tree should contain Database")
	}

	if !strings.Contains(result, "Host") {
		t.Error("tree should contain Host (child of Database)")
	}
}

func TestRenderTree_WithFocus(t *testing.T) {
	nodes := []ast.Node{
		{
			Name: "Database",
			Kind: ast.KindStruct,
			Children: []ast.Node{
				{Name: "Host", Kind: ast.KindPrimitive},
			},
		},
	}

	focusPath := path.Root().Child("Database")
	html := RenderTree(nodes, focusPath, nil)
	result := renderToString(html)

	if !strings.Contains(result, "tree-node--active") {
		t.Error("tree should highlight focused node")
	}
}

func TestRenderTree_WithSlice(t *testing.T) {
	nodes := []ast.Node{
		{
			Name:        "Services",
			Kind:        ast.KindSlice,
			ElementKind: ast.KindStruct,
		},
	}

	values := map[string]any{
		"Services": []map[string]any{
			{"Name": "API"},
			{"Name": "Worker"},
		},
	}

	html := RenderTree(nodes, path.Root(), values)
	result := renderToString(html)

	if !strings.Contains(result, "Services") {
		t.Error("tree should contain Services")
	}
}
