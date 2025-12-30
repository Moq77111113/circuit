package layout

import (
	"strings"
	"testing"
)

func TestRenderActions_Empty(t *testing.T) {
	result := renderActions([]ActionButton{})
	if result != nil {
		t.Fatal("expected nil for empty actions")
	}
}

func TestRenderActions_SingleAction(t *testing.T) {
	actions := []ActionButton{
		{
			Name:                "test",
			Label:               "Test Action",
			Description:         "Test description",
			RequireConfirmation: false,
		},
	}

	result := renderActions(actions)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	html := renderToString(result)
	if !strings.Contains(html, "Test Action") {
		t.Error("expected action label in output")
	}
	if !strings.Contains(html, "execute:test") {
		t.Error("expected execute:test value in output")
	}
}

func TestRenderActions_MultipleActions(t *testing.T) {
	actions := []ActionButton{
		{Name: "action1", Label: "Action 1"},
		{Name: "action2", Label: "Action 2"},
		{Name: "action3", Label: "Action 3"},
	}

	result := renderActions(actions)
	html := renderToString(result)

	if !strings.Contains(html, "Action 1") {
		t.Error("expected Action 1 in output")
	}
	if !strings.Contains(html, "Action 2") {
		t.Error("expected Action 2 in output")
	}
	if !strings.Contains(html, "Action 3") {
		t.Error("expected Action 3 in output")
	}
}

func TestRenderActions_WithConfirmation(t *testing.T) {
	actions := []ActionButton{
		{
			Name:                "dangerous",
			Label:               "Dangerous Action",
			RequireConfirmation: true,
		},
	}

	result := renderActions(actions)
	html := renderToString(result)

	if !strings.Contains(html, "onclick") {
		t.Error("expected onclick attribute for confirmation")
	}
	if !strings.Contains(html, "confirmAction") {
		t.Error("expected confirmAction function call")
	}
}

func TestRenderActions_WithoutConfirmation(t *testing.T) {
	actions := []ActionButton{
		{
			Name:                "safe",
			Label:               "Safe Action",
			RequireConfirmation: false,
		},
	}

	result := renderActions(actions)
	html := renderToString(result)

	if strings.Contains(html, "onclick") {
		t.Error("expected no onclick attribute without confirmation")
	}
}
