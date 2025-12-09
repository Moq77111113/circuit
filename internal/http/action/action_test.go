package action

import (
	"net/url"
	"testing"
)

func TestParseAction_Save(t *testing.T) {
	form := url.Values{
		"action": {"save"},
	}

	action := Parse(form)

	if action.Type != ActionSave {
		t.Errorf("expected action type %s, got %s", ActionSave, action.Type)
	}
}

func TestParseAction_Add(t *testing.T) {
	form := url.Values{
		"action": {"add:tags"},
	}

	action := Parse(form)

	if action.Type != ActionAdd {
		t.Errorf("expected action type %s, got %s", ActionAdd, action.Type)
	}
	if action.Field != "tags" {
		t.Errorf("expected field tags, got %s", action.Field)
	}
}

func TestParseAction_Remove(t *testing.T) {
	form := url.Values{
		"action": {"remove:items:2"},
	}

	action := Parse(form)

	if action.Type != ActionRemove {
		t.Errorf("expected action type %s, got %s", ActionRemove, action.Type)
	}
	if action.Field != "items" {
		t.Errorf("expected field items, got %s", action.Field)
	}
	if action.Index != 2 {
		t.Errorf("expected index 2, got %d", action.Index)
	}
}

func TestParseAction_InvalidIndex(t *testing.T) {
	form := url.Values{
		"action": {"remove:items:invalid"},
	}

	action := Parse(form)

	if action.Type != ActionSave {
		t.Errorf("expected fallback to save action, got %s", action.Type)
	}
}

func TestParseAction_Empty(t *testing.T) {
	form := url.Values{}

	action := Parse(form)

	if action.Type != ActionSave {
		t.Errorf("expected default save action, got %s", action.Type)
	}
}

func TestParseAction_Malformed(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"add without field", "add"},
		{"add with empty field", "add:"},
		{"remove without field", "remove"},
		{"remove without index", "remove:items"},
		{"remove with empty index", "remove:items:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{
				"action": {tt.value},
			}

			action := Parse(form)

			if action.Type != ActionSave {
				t.Errorf("expected fallback to save for %q, got %s", tt.value, action.Type)
			}
		})
	}
}
