package collapsible

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ui/styles"
)

func renderToString(node g.Node) string {
	var sb strings.Builder
	_ = node.Render(&sb)
	return sb.String()
}

func TestCollapsible(t *testing.T) {
	cfg := Config{
		ID:    "test-collapsible",
		Title: "Test Section",
		Depth: 0,
		Count: 3,
	}

	result := Collapsible(cfg, []g.Node{h.Div(g.Text("child"))})
	html := renderToString(result)

	if !strings.Contains(html, styles.Collapsible) {
		t.Error("should contain collapsible class")
	}
	if !strings.Contains(html, "id=\"test-collapsible\"") {
		t.Error("should contain ID attribute")
	}
	if !strings.Contains(html, "Test Section") {
		t.Error("should contain title")
	}
	if !strings.Contains(html, "(3)") {
		t.Error("should contain count badge")
	}
}

func TestCollapsibleDepth(t *testing.T) {
	tests := []struct {
		depth int
		want  string
	}{
		{0, styles.CollapsibleDepth0},
		{1, styles.CollapsibleDepth1},
		{2, styles.CollapsibleDepth2},
		{3, styles.CollapsibleDepth3},
		{4, styles.CollapsibleDepth4},
	}

	for _, tt := range tests {
		cfg := Config{Title: "Test", Depth: tt.depth}
		result := Collapsible(cfg, nil)
		html := renderToString(result)

		if !strings.Contains(html, tt.want) {
			t.Errorf("depth %d should contain class %q", tt.depth, tt.want)
		}
	}
}

func TestCollapsibleCollapsed(t *testing.T) {
	cfg := Config{Title: "Test", Collapsed: true}
	result := Collapsible(cfg, nil)
	html := renderToString(result)

	if !strings.Contains(html, styles.CollapsibleCollapsed) {
		t.Error("collapsed state should add collapsed class")
	}
}

func TestHeader(t *testing.T) {
	result := Header("My Title", 5)
	html := renderToString(result)

	if !strings.Contains(html, styles.CollapsibleHeader) {
		t.Error("should contain header class")
	}
	if !strings.Contains(html, styles.IconChevronDown) {
		t.Error("should contain chevron icon class")
	}
	if !strings.Contains(html, "My Title") {
		t.Error("should contain title text")
	}
	if !strings.Contains(html, "(5)") {
		t.Error("should contain count")
	}
}

func TestHeaderNoCount(t *testing.T) {
	result := Header("Title", 0)
	html := renderToString(result)

	if strings.Contains(html, "(0)") {
		t.Error("should not show count when zero")
	}
}

func TestBody(t *testing.T) {
	result := Body([]g.Node{h.P(g.Text("content"))})
	html := renderToString(result)

	if !strings.Contains(html, styles.CollapsibleBody) {
		t.Error("should contain body class")
	}
	if !strings.Contains(html, "content") {
		t.Error("should contain child content")
	}
}
