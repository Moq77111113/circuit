package containers

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/tags"
)

func renderToString(node any) string {
	var buf bytes.Buffer
	if n, ok := node.(interface{ Render(io.Writer) error }); ok {
		_ = n.Render(&buf)
	}
	return buf.String()
}

func TestSlice_EmptyStringSlice(t *testing.T) {
	field := tags.Field{
		Name:        "Tags",
		IsSlice:     true,
		ElementType: "string",
		InputType:   tags.TypeText,
	}

	node := Slice(field, []string{})
	html := renderToString(node)

	if !strings.Contains(html, "Add Item") {
		t.Error("expected 'Add Item' button for empty slice")
	}
}

func TestSlice_StringSliceWithValues(t *testing.T) {
	field := tags.Field{
		Name:        "Tags",
		IsSlice:     true,
		ElementType: "string",
		InputType:   tags.TypeText,
	}

	node := Slice(field, []string{"go", "rust"})
	html := renderToString(node)

	if !strings.Contains(html, "Tags.0") {
		t.Error("expected indexed field name Tags.0")
	}
	if !strings.Contains(html, "Tags.1") {
		t.Error("expected indexed field name Tags.1")
	}
	if !strings.Contains(html, "value=\"go\"") {
		t.Error("expected value 'go'")
	}
	if !strings.Contains(html, "value=\"rust\"") {
		t.Error("expected value 'rust'")
	}
}

func TestSlice_IntSlice(t *testing.T) {
	field := tags.Field{
		Name:        "Ports",
		IsSlice:     true,
		ElementType: "int",
		InputType:   tags.TypeNumber,
	}

	node := Slice(field, []int{8080, 9090})
	html := renderToString(node)

	if !strings.Contains(html, "Ports.0") {
		t.Error("expected indexed field name Ports.0")
	}
	if !strings.Contains(html, "value=\"8080\"") {
		t.Error("expected value '8080'")
	}
	if !strings.Contains(html, "type=\"number\"") {
		t.Error("expected number input type")
	}
}

func TestSlice_RemoveButtons(t *testing.T) {
	field := tags.Field{
		Name:        "Tags",
		IsSlice:     true,
		ElementType: "string",
		InputType:   tags.TypeText,
	}

	node := Slice(field, []string{"go", "rust"})
	html := renderToString(node)

	if !strings.Contains(html, "remove:Tags:0") {
		t.Error("expected remove button for index 0")
	}
	if !strings.Contains(html, "remove:Tags:1") {
		t.Error("expected remove button for index 1")
	}
}

func TestSlice_AddButton(t *testing.T) {
	field := tags.Field{
		Name:        "Tags",
		IsSlice:     true,
		ElementType: "string",
		InputType:   tags.TypeText,
	}

	node := Slice(field, []string{"go"})
	html := renderToString(node)

	if !strings.Contains(html, "add:Tags") {
		t.Error("expected add button with correct action")
	}
}

func TestSliceWithDepth_ShallowFull(t *testing.T) {
	field := tags.Field{
		Name: "Services",
		Fields: []tags.Field{
			{Name: "Name", Type: "string"},
		},
		IsSlice:     true,
		ElementType: "struct",
	}
	value := []struct{ Name string }{{Name: "API"}}


	for _, depth := range []int{0, 1} {
		node := SliceWithDepth(field, value, depth)
		html := renderToString(node)

		if strings.Contains(html, "slice__item--compact") {
			t.Errorf("depth %d should NOT use compact rendering", depth)
		}
	}
}

func TestSliceWithDepth_DeepCompact(t *testing.T) {
	field := tags.Field{
		Name: "Services",
		Fields: []tags.Field{
			{Name: "Name", Type: "string"},
		},
		IsSlice:     true,
		ElementType: "struct",
	}
	value := []struct{ Name string }{{Name: "API"}}

	node := SliceWithDepth(field, value, 2)
	html := renderToString(node)

	if !strings.Contains(html, "slice__item--compact") {
		t.Error("depth 2 should use compact rendering")
	}
}

func TestSliceWithDepth_Summary(t *testing.T) {
	field := tags.Field{
		Name: "Services",
		Fields: []tags.Field{
			{Name: "Name", Type: "string"},
			{Name: "Type", Type: "string"},
		},
		IsSlice:     true,
		ElementType: "struct",
	}
	value := []struct {
		Name string
		Type string
	}{{Name: "User Service", Type: "HTTP"}}

	node := SliceWithDepth(field, value, 1)
	html := renderToString(node)

	if !strings.Contains(html, "User Service") || !strings.Contains(html, "HTTP") {
		t.Error("should include summary in header for struct slices")
	}
}

func TestSliceWithDepth_PrimitiveNoSummary(t *testing.T) {
	field := tags.Field{
		Name:        "Tags",
		IsSlice:     true,
		ElementType: "string",
		InputType:   tags.TypeText,
	}
	value := []string{"go", "web"}

	node := SliceWithDepth(field, value, 1)
	html := renderToString(node)

	if !strings.Contains(html, "slice__header") {
		t.Error("should render header for primitive slice")
	}
}
