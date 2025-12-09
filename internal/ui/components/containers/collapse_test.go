package containers

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func TestCollapsibleHeader(t *testing.T) {
	node := CollapsibleHeader("Test Title", 5, false)
	html := renderNode(node)

	if !strings.Contains(html, "Test Title") {
		t.Error("Header should contain title")
	}
	if !strings.Contains(html, "(5)") {
		t.Error("Header should contain count")
	}
	if !strings.Contains(html, "slice__header") {
		t.Error("Header should have slice__header class")
	}
	if !strings.Contains(html, "onclick=\"toggleCollapse(this)\"") {
		t.Error("Header should have onclick handler")
	}
}

func TestCollapsibleContainer(t *testing.T) {
	tests := []struct {
		depth          int
		shouldCollapse bool
	}{
		{0, false},
		{1, false},
		{2, true},
	}

	for _, tt := range tests {
		node := CollapsibleContainer(tt.depth, h.Div())
		html := renderNode(node)

		if !strings.Contains(html, DepthClass(tt.depth)) {
			t.Errorf("Container should have depth class for depth %d", tt.depth)
		}

		if tt.shouldCollapse {
			if !strings.Contains(html, "collapsed") {
				t.Errorf("Container should have collapsed class for depth %d", tt.depth)
			}
		} else {
			if strings.Contains(html, "collapsed") {
				t.Errorf("Container should NOT have collapsed class for depth %d", tt.depth)
			}
		}
	}
}

func TestSectionHeader(t *testing.T) {
	t.Run("Collapsible", func(t *testing.T) {
		node := SectionHeader("Title", true)
		html := renderNode(node)
		if !strings.Contains(html, "onclick=\"toggleCollapse(this)\"") {
			t.Error("Collapsible section header should have onclick")
		}
	})

	t.Run("Not Collapsible", func(t *testing.T) {
		node := SectionHeader("Title", false)
		html := renderNode(node)
		if strings.Contains(html, "onclick=\"toggleCollapse(this)\"") {
			t.Error("Non-collapsible section header should NOT have onclick")
		}
	})
}

func renderNode(n g.Node) string {
	var b strings.Builder
	_ = n.Render(&b)
	return b.String()
}
