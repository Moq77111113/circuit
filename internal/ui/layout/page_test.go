package layout

import (
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/tags"
	g "maragu.dev/gomponents"
)

func renderToString(node g.Node) string {
	var sb strings.Builder
	_ = node.Render(&sb)
	return sb.String()
}

func TestPage_Complete(t *testing.T) {
	s := ast.Schema{
		Name: "Config",
		Nodes: ast.FromTags([]tags.Field{
			{Name: "Host", Type: "string", InputType: "text", Help: "Server host"},
			{Name: "Port", Type: "int", InputType: "number", Help: "Server port"},
		}),
	}

	html := renderToString(Page(s, nil, path.Root(), PageOptions{Brand: true}))

	// Check HTML structure
	if !strings.Contains(html, "<!doctype html>") && !strings.Contains(html, "<html") {
		t.Error("expected html element")
	}
	if !strings.Contains(html, "<head>") {
		t.Error("expected head element")
	}
	if !strings.Contains(html, "<style>") {
		t.Error("expected embedded style")
	}
	if !strings.Contains(html, "<body>") {
		t.Error("expected body element")
	}
	if !strings.Contains(html, "Config Configuration") {
		t.Error("expected default title")
	}
}

func TestPage_CustomTitle(t *testing.T) {
	s := ast.Schema{
		Name:  "Config",
		Nodes: ast.FromTags([]tags.Field{}),
	}

	html := renderToString(Page(s, nil, path.Root(), PageOptions{Title: "My Custom Settings", Brand: true}))

	if !strings.Contains(html, "My Custom Settings") {
		t.Error("expected custom title")
	}
}

func TestPage_WithValues(t *testing.T) {
	s := ast.Schema{
		Name: "Config",
		Nodes: ast.FromTags([]tags.Field{
			{Name: "Host", Type: "string", InputType: "text"},
		}),
	}

	values := map[string]any{
		"Host": "example.com",
	}

	html := renderToString(Page(s, values, path.Root(), PageOptions{Brand: true}))

	if !strings.Contains(html, `value="example.com"`) {
		t.Error("expected value in rendered page")
	}
}

func TestPage_ContainsCSS(t *testing.T) {
	s := ast.Schema{
		Name:  "Config",
		Nodes: ast.FromTags([]tags.Field{}),
	}

	html := renderToString(Page(s, nil, path.Root(), PageOptions{Brand: true}))

	// Check for key CSS classes
	if !strings.Contains(html, ".app__container") {
		t.Error("expected app__container CSS")
	}
	if !strings.Contains(html, ".form") {
		t.Error("expected form CSS")
	}
	if !strings.Contains(html, ".field__input") {
		t.Error("expected field__input CSS")
	}
}

func TestPage_ContainsBranding(t *testing.T) {
	s := ast.Schema{
		Name:  "Config",
		Nodes: ast.FromTags([]tags.Field{}),
	}

	html := renderToString(Page(s, nil, path.Root(), PageOptions{Brand: true}))

	if !strings.Contains(html, "Powered by") {
		t.Error("expected 'Powered by' text")
	}
	if !strings.Contains(html, "Circuit") {
		t.Error("expected Circuit branding")
	}
	if !strings.Contains(html, "github.com/moq77111113/circuit") {
		t.Error("expected GitHub link")
	}
	if !strings.Contains(html, "<footer") {
		t.Error("expected footer class")
	}
}

func TestPage_WithoutBranding(t *testing.T) {
	s := ast.Schema{
		Name:  "Config",
		Nodes: ast.FromTags([]tags.Field{}),
	}

	html := renderToString(Page(s, nil, path.Root(), PageOptions{}))

	if strings.Contains(html, "Powered by") {
		t.Error("did not expect 'Powered by' text")
	}
	if strings.Contains(html, "Circuit") {
		t.Error("did not expect Circuit branding")
	}
	if strings.Contains(html, "github.com/moq77111113/circuit") {
		t.Error("did not expect GitHub link")
	}
	if strings.Contains(html, "<footer") {
		t.Error("did not expect footer class")
	}
}
