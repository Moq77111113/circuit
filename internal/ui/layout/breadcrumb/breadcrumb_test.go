package breadcrumb

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

func TestRenderBreadcrumb_Root(t *testing.T) {
	nodes := []ast.Node{
		{Name: "AppName", Kind: ast.KindPrimitive},
	}

	html := RenderBreadcrumb(path.Root(), nodes)
	result := renderToString(html)

	if !strings.Contains(result, "Config") {
		t.Error("breadcrumb should show 'Config' for root")
	}

	if strings.Contains(result, "breadcrumb__separator") {
		t.Error("breadcrumb should have no separators at root")
	}
}

func TestRenderBreadcrumb_OneLevelDeep(t *testing.T) {
	nodes := []ast.Node{
		{Name: "Database", Kind: ast.KindStruct},
	}

	currentPath := path.Root().Child("Database")
	html := RenderBreadcrumb(currentPath, nodes)
	result := renderToString(html)

	if !strings.Contains(result, "Config") {
		t.Error("breadcrumb should contain Config")
	}

	if !strings.Contains(result, "Database") {
		t.Error("breadcrumb should contain Database")
	}

	if !strings.Contains(result, ">") {
		t.Error("breadcrumb should have separator")
	}
}

func TestRenderBreadcrumb_DeepNesting(t *testing.T) {
	nodes := []ast.Node{
		{
			Name: "Server",
			Kind: ast.KindStruct,
			Children: []ast.Node{
				{Name: "RateLimit", Kind: ast.KindStruct},
			},
		},
	}

	currentPath := path.Root().Child("Server").Child("RateLimit")
	html := RenderBreadcrumb(currentPath, nodes)
	result := renderToString(html)

	if !strings.Contains(result, "Config") {
		t.Error("breadcrumb should contain Config")
	}

	if !strings.Contains(result, "Server") {
		t.Error("breadcrumb should contain Server")
	}

	if !strings.Contains(result, "RateLimit") {
		t.Error("breadcrumb should contain RateLimit")
	}
}

func TestRenderBreadcrumb_SliceItem(t *testing.T) {
	nodes := []ast.Node{
		{
			Name:        "Services",
			Kind:        ast.KindSlice,
			ElementKind: ast.KindStruct,
		},
	}

	currentPath := path.Root().Child("Services").Index(0)
	html := RenderBreadcrumb(currentPath, nodes)
	result := renderToString(html)

	if !strings.Contains(result, "Services") {
		t.Error("breadcrumb should contain Services")
	}

	if !strings.Contains(result, "[0]") {
		t.Errorf("breadcrumb should contain slice index [0], got: %s", result)
	}
}

func renderToString(node g.Node) string {
	var sb strings.Builder
	_ = node.Render(&sb)
	return sb.String()
}
