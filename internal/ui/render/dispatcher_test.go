package render

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/tags"
)

// renderToString is a helper to convert gomponents to string for testing
func renderToString(node g.Node) string {
	var sb strings.Builder
	_ = node.Render(&sb)
	return sb.String()
}

func TestDispatcher_RoutesToPrimitive(t *testing.T) {
	d := NewDispatcher()
	node := schema.Node{
		Name:      "Port",
		Kind:      schema.KindPrimitive,
		ValueType: schema.ValueInt,
		InputType: tags.TypeNumber,
	}
	ctx := Context{
		Path:  schema.NewPath("Port"),
		Value: 8080,
		Depth: 0,
	}

	html := renderToString(d.Render(node, ctx))

	if !strings.Contains(html, `class="field"`) {
		t.Error("expected field wrapper class")
	}

	if !strings.Contains(html, `name="Port"`) {
		t.Error("expected name attribute")
	}
}

func TestDispatcher_RoutesToStruct(t *testing.T) {
	d := NewDispatcher()
	node := schema.Node{
		Name: "Database",
		Kind: schema.KindStruct,
		Children: []schema.Node{
			{
				Name:      "Host",
				Kind:      schema.KindPrimitive,
				ValueType: schema.ValueString,
				InputType: tags.TypeText,
			},
		},
	}
	ctx := Context{
		Path: schema.NewPath("Database"),
		Value: struct {
			Host string
		}{Host: "localhost"},
		Depth: 0,
	}

	html := renderToString(d.Render(node, ctx))

	if !strings.Contains(html, `<section`) {
		t.Error("expected section element")
	}
		if !strings.Contains(html, `name="Database.Host"`) {
		t.Error("expected nested field path")
	}
}

func TestDispatcher_RoutesToSlice(t *testing.T) {
	d := NewDispatcher()
	node := schema.Node{
		Name:        "Tags",
		Kind:        schema.KindSlice,
		ElementKind: schema.KindPrimitive,
		ValueType:   schema.ValueString,
		InputType:   tags.TypeText,
	}
	ctx := Context{
		Path:  schema.NewPath("Tags"),
		Value: []string{"go", "web"},
		Depth: 0,
	}

	html := renderToString(d.Render(node, ctx))

	if !strings.Contains(html, `name="Tags.0"`) {
		t.Error("expected first indexed field")
	}
	if !strings.Contains(html, `value="go"`) {
		t.Error("expected first value")
	}
}

func TestContext_ChildPath(t *testing.T) {
	ctx := Context{
		Path:  schema.NewPath("Database"),
		Value: nil,
		Depth: 0,
	}

	childCtx := Context{
		Path:  ctx.Path.Child("Host"),
		Value: "localhost",
		Depth: ctx.Depth + 1,
	}

	got := childCtx.Path.String()
	want := "Database.Host"

	if got != want {
		t.Errorf("child path = %q, want %q", got, want)
	}

	if childCtx.Depth != 1 {
		t.Errorf("child depth = %d, want 1", childCtx.Depth)
	}
}
