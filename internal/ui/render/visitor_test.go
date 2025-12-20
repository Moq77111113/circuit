package render

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/tags"
)

// renderToString is a helper to convert gomponents to string for testing
func renderToString(node g.Node) string {
	var sb strings.Builder
	_ = node.Render(&sb)
	return sb.String()
}

func TestRenderVisitor_Primitive(t *testing.T) {
	nodes := []ast.Node{
		{
			Name:      "Port",
			Kind:      ast.KindPrimitive,
			ValueType: ast.ValueInt,
			UI:        &ast.UIMetadata{InputType: tags.TypeNumber},
		},
	}
	values := map[string]any{
		"Port": 8080,
	}

	html := renderToString(Render(nodes, values))

	if !strings.Contains(html, `class="field"`) {
		t.Error("expected field wrapper class")
	}
	if !strings.Contains(html, `name="Port"`) {
		t.Error("expected name attribute")
	}
	if !strings.Contains(html, `value="8080"`) {
		t.Error("expected value 8080")
	}
}

func TestRenderVisitor_PrimitiveString(t *testing.T) {
	nodes := []ast.Node{
		{
			Name:      "Name",
			Kind:      ast.KindPrimitive,
			ValueType: ast.ValueString,
			UI:        &ast.UIMetadata{InputType: tags.TypeText},
		},
	}
	values := map[string]any{
		"Name": "test-server",
	}

	html := renderToString(Render(nodes, values))

	if !strings.Contains(html, `type="text"`) {
		t.Error("expected text input type")
	}
	if !strings.Contains(html, `value="test-server"`) {
		t.Error("expected value test-server")
	}
}

func TestRenderVisitor_Struct(t *testing.T) {
	nodes := []ast.Node{
		{
			Name: "Database",
			Kind: ast.KindStruct,
			Children: []ast.Node{
				{
					Name:      "Host",
					Kind:      ast.KindPrimitive,
					ValueType: ast.ValueString,
					UI:        &ast.UIMetadata{InputType: tags.TypeText},
				},
				{
					Name:      "Port",
					Kind:      ast.KindPrimitive,
					ValueType: ast.ValueInt,
					UI:        &ast.UIMetadata{InputType: tags.TypeNumber},
				},
			},
		},
	}
	values := map[string]any{
		"Database.Host": "localhost",
		"Database.Port": 5432,
	}

	html := renderToString(Render(nodes, values))

	// Check nested paths
	if !strings.Contains(html, `name="Database.Host"`) {
		t.Error("expected nested field path Database.Host")
	}
	if !strings.Contains(html, `name="Database.Port"`) {
		t.Error("expected nested field path Database.Port")
	}
	if !strings.Contains(html, `value="localhost"`) {
		t.Error("expected Host value")
	}
	if !strings.Contains(html, `value="5432"`) {
		t.Error("expected Port value")
	}
}

func TestRenderVisitor_SlicePrimitive(t *testing.T) {
	nodes := []ast.Node{
		{
			Name:        "Tags",
			Kind:        ast.KindSlice,
			ElementKind: ast.KindPrimitive,
			ValueType:   ast.ValueString,
			UI:          &ast.UIMetadata{InputType: tags.TypeText},
		},
	}
	values := map[string]any{
		"Tags": []string{"go", "web", "api"},
	}

	html := renderToString(Render(nodes, values))

	// Check indexed paths
	if !strings.Contains(html, `name="Tags.0"`) {
		t.Error("expected Tags.0 field")
	}
	if !strings.Contains(html, `name="Tags.1"`) {
		t.Error("expected Tags.1 field")
	}
	if !strings.Contains(html, `name="Tags.2"`) {
		t.Error("expected Tags.2 field")
	}

	// Check values
	if !strings.Contains(html, `value="go"`) {
		t.Error("expected value go")
	}
	if !strings.Contains(html, `value="web"`) {
		t.Error("expected value web")
	}
	if !strings.Contains(html, `value="api"`) {
		t.Error("expected value api")
	}

	// Check add button
	if !strings.Contains(html, `name="action"`) {
		t.Error("expected add button")
	}
	if !strings.Contains(html, `value="add_Tags"`) {
		t.Error("expected add_Tags action")
	}
}

func TestRenderVisitor_SliceStruct(t *testing.T) {
	nodes := []ast.Node{
		{
			Name:        "Services",
			Kind:        ast.KindSlice,
			ElementKind: ast.KindStruct,
			Children: []ast.Node{
				{
					Name:      "Name",
					Kind:      ast.KindPrimitive,
					ValueType: ast.ValueString,
					UI:        &ast.UIMetadata{InputType: tags.TypeText},
				},
				{
					Name:      "Port",
					Kind:      ast.KindPrimitive,
					ValueType: ast.ValueInt,
					UI:        &ast.UIMetadata{InputType: tags.TypeNumber},
				},
			},
		},
	}
	values := map[string]any{
		"Services": []struct {
			Name string
			Port int
		}{
			{Name: "api", Port: 8080},
			{Name: "web", Port: 3000},
		},
	}

	html := renderToString(Render(nodes, values))

	// Check struct slice items (summary should be rendered)
	if !strings.Contains(html, "api") {
		t.Error("expected service name in summary")
	}
	if !strings.Contains(html, "web") {
		t.Error("expected second service name in summary")
	}

	// Check collapsible container for slice (has depth class)
	if !strings.Contains(html, `class="slice slice--depth-0"`) &&
		!strings.Contains(html, `class="slice slice--depth-0 collapsed"`) {
		t.Error("expected slice container with depth class")
	}
}

func TestRenderVisitor_NestedStruct(t *testing.T) {
	nodes := []ast.Node{
		{
			Name: "Server",
			Kind: ast.KindStruct,
			Children: []ast.Node{
				{
					Name: "Database",
					Kind: ast.KindStruct,
					Children: []ast.Node{
						{
							Name:      "Host",
							Kind:      ast.KindPrimitive,
							ValueType: ast.ValueString,
							UI:        &ast.UIMetadata{InputType: tags.TypeText},
						},
					},
				},
			},
		},
	}
	values := map[string]any{
		"Server.Database.Host": "localhost",
	}

	html := renderToString(Render(nodes, values))

	// Check deeply nested path
	if !strings.Contains(html, `name="Server.Database.Host"`) {
		t.Error("expected deeply nested path Server.Database.Host")
	}
	if !strings.Contains(html, `value="localhost"`) {
		t.Error("expected value localhost")
	}
}

func TestRenderVisitor_EmptySlice(t *testing.T) {
	nodes := []ast.Node{
		{
			Name:        "Tags",
			Kind:        ast.KindSlice,
			ElementKind: ast.KindPrimitive,
			ValueType:   ast.ValueString,
			UI:          &ast.UIMetadata{InputType: tags.TypeText},
		},
	}
	values := map[string]any{
		"Tags": []string{},
	}

	html := renderToString(Render(nodes, values))

	// Empty slice should still have add button
	if !strings.Contains(html, `value="add_Tags"`) {
		t.Error("expected add button for empty slice")
	}

	// Should not have any indexed fields
	if strings.Contains(html, `name="Tags.0"`) {
		t.Error("should not have Tags.0 for empty slice")
	}
}
