package containers

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestSection_CollapsibleUsesContainer(t *testing.T) {
	nodes := []g.Node{g.Text("content")}
	section := Section("Server", nodes, true, false)
	html := renderNode(section)

	if !strings.Contains(html, "container--depth-0") {
		t.Error("Collapsible section should use container--depth-0 class")
	}
	if !strings.Contains(html, "container__header") {
		t.Error("Collapsible section should use container__header")
	}
	if !strings.Contains(html, "container__body") {
		t.Error("Collapsible section should use container__body")
	}
	if !strings.Contains(html, "<section") {
		t.Error("Collapsible section should still use semantic <section> tag")
	}
}

func TestSection_CollapsibleNoCount(t *testing.T) {
	nodes := []g.Node{g.Text("content")}
	section := Section("Database", nodes, true, false)
	html := renderNode(section)

	if strings.Contains(html, "container__count") {
		t.Error("Section should not display item count")
	}
	if strings.Contains(html, "slice__count") {
		t.Error("Section should not use slice__count class")
	}
}

func TestSection_NonCollapsibleKeepsSemantic(t *testing.T) {
	nodes := []g.Node{g.Text("content")}
	section := Section("Server", nodes, false, false)
	html := renderNode(section)

	if !strings.Contains(html, "<section") {
		t.Error("Non-collapsible section should use semantic <section> tag")
	}
	if !strings.Contains(html, "section__title") {
		t.Error("Non-collapsible section should have section__title")
	}
}

func TestSection_CollapsedState(t *testing.T) {
	nodes := []g.Node{g.Text("content")}
	section := Section("Pool", nodes, true, true)
	html := renderNode(section)

	if !strings.Contains(html, "collapsed") {
		t.Error("Collapsed section should have 'collapsed' class")
	}
}

func TestSection_MatchesSliceStyle(t *testing.T) {
	nodes := []g.Node{g.Text("test content")}
	section := Section("Database", nodes, true, false)
	sectionHTML := renderNode(section)

	if !strings.Contains(sectionHTML, "container--depth-0") {
		t.Error("Section should use container--depth-0 class")
	}
	if !strings.Contains(sectionHTML, "container__header") {
		t.Error("Section should use container__header class")
	}
	if !strings.Contains(sectionHTML, "container__body") {
		t.Error("Section should use container__body class")
	}
	if !strings.Contains(sectionHTML, "container__chevron") {
		t.Error("Section should use container__chevron class")
	}
	if !strings.Contains(sectionHTML, "container__title") {
		t.Error("Section should use container__title class")
	}
	if !strings.Contains(sectionHTML, "<section") {
		t.Error("Section should use semantic <section> HTML tag")
	}
}

func TestSection_IndexedTitleCleanup(t *testing.T) {
	nodes := []g.Node{g.Text("auth config")}
	section := Section("Services.0.Auth", nodes, true, false)
	html := renderNode(section)

	if !strings.Contains(html, ">Auth<") {
		t.Errorf("Section should display clean title 'Auth', got:\n%s", html)
	}
	if strings.Contains(html, "container__title\">Services.0.Auth<") {
		t.Errorf("Section title text should not display full indexed path, got:\n%s", html)
	}
}
