package containers

import (
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/tags"
)

func TestRender_OneField(t *testing.T) {
	field := tags.Field{
		Name: "Services",
		Fields: []tags.Field{
			{Name: "Name", Type: "string"},
		},
	}
	value := struct{ Name string }{Name: "User Service"}

	node := Render(field, 0, value, 2)
	html := renderNode(node)

	if !strings.Contains(html, "Name: User Service") {
		t.Error("Should contain field value")
	}
	if !strings.Contains(html, "slice__item--compact") {
		t.Error("Should have compact class")
	}
}

func TestRender_ThreeFields(t *testing.T) {
	field := tags.Field{
		Name: "Services",
		Fields: []tags.Field{
			{Name: "Name", Type: "string"},
			{Name: "Type", Type: "string"},
			{Name: "Enabled", Type: "bool"},
		},
	}
	value := struct {
		Name    string
		Type    string
		Enabled bool
	}{Name: "API", Type: "HTTP", Enabled: true}

	node := Render(field, 0, value, 2)
	html := renderNode(node)

	if !strings.Contains(html, "Name: API") {
		t.Error("Should contain Name field")
	}
	if !strings.Contains(html, "Type: HTTP") {
		t.Error("Should contain Type field")
	}
	if !strings.Contains(html, "Enabled: true") {
		t.Error("Should contain Enabled field")
	}
}

func TestRender_IndexNumber(t *testing.T) {
	field := tags.Field{
		Name:   "Services",
		Fields: []tags.Field{{Name: "Name", Type: "string"}},
	}
	value := struct{ Name string }{Name: "Test"}

	node := Render(field, 2, value, 2)
	html := renderNode(node)

	if !strings.Contains(html, "#3") {
		t.Error("Should display index number (0-based + 1)")
	}
}

func TestRender_RemoveButton(t *testing.T) {
	field := tags.Field{
		Name:   "Services",
		Fields: []tags.Field{{Name: "Name", Type: "string"}},
	}
	value := struct{ Name string }{Name: "Test"}

	node := Render(field, 0, value, 2)
	html := renderNode(node)

	if !strings.Contains(html, "Remove") {
		t.Error("Should have Remove button")
	}
	if !strings.Contains(html, "remove:Services:0") {
		t.Error("Should have correct remove action")
	}
}

func TestRender_Chevron(t *testing.T) {
	field := tags.Field{
		Name:   "Services",
		Fields: []tags.Field{{Name: "Name", Type: "string"}},
	}
	value := struct{ Name string }{Name: "Test"}

	node := Render(field, 0, value, 2)
	html := renderNode(node)

	if !strings.Contains(html, "▼") {
		t.Error("Should have chevron indicator")
	}
}

func TestRender_CSSClass(t *testing.T) {
	field := tags.Field{
		Name:   "Services",
		Fields: []tags.Field{{Name: "Name", Type: "string"}},
	}
	value := struct{ Name string }{Name: "Test"}

	node := Render(field, 0, value, 2)
	html := renderNode(node)

	if !strings.Contains(html, "slice__item--compact") {
		t.Error("Should have slice__item--compact class")
	}
}
