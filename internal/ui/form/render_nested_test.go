package form

import (
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/tags"
)

func TestRenderNested_ThreeLevels(t *testing.T) {
	field := tags.Field{
		Name: "Services",
		Fields: []tags.Field{
			{Name: "Name", Type: "string"},
			{Name: "Enabled", Type: "bool"},
			{Name: "Endpoints", IsSlice: true, ElementType: "struct", Fields: []tags.Field{
				{Name: "Path", Type: "string"},
				{Name: "Method", Type: "string"},
				{Name: "AllowedRoles", IsSlice: true, ElementType: "string"},
			}},
		},
		IsSlice:     true,
		ElementType: "struct",
	}

	value := []struct {
		Name      string
		Enabled   bool
		Endpoints []struct {
			Path         string
			Method       string
			AllowedRoles []string
		}
	}{
		{
			Name:    "User Service",
			Enabled: true,
			Endpoints: []struct {
				Path         string
				Method       string
				AllowedRoles []string
			}{
				{Path: "/api/v1/users", Method: "GET", AllowedRoles: []string{"admin", "user"}},
			},
		},
	}

	node := renderInput(field, value)
	html := renderToString(node)

	if !strings.Contains(html, "Services") {
		t.Error("Should render Services slice")
	}
	if !strings.Contains(html, "Endpoints") {
		t.Error("Should render Endpoints slice")
	}
	if !strings.Contains(html, "AllowedRoles") {
		t.Error("Should render AllowedRoles slice")
	}
}

func TestRenderNested_DepthExpansion(t *testing.T) {
	field := tags.Field{
		Name: "Services",
		Fields: []tags.Field{
			{Name: "Name", Type: "string"},
			{Name: "Endpoints", IsSlice: true, ElementType: "struct", Fields: []tags.Field{
				{Name: "Path", Type: "string"},
			}},
		},
		IsSlice:     true,
		ElementType: "struct",
	}

	value := []struct {
		Name      string
		Endpoints []struct{ Path string }
	}{
		{Name: "Test", Endpoints: []struct{ Path string }{{Path: "/api"}}},
	}

	html := renderToString(renderInput(field, value))

	if strings.Contains(html, "slice--depth-0 collapsed") {
		t.Error("Depth 0 should be expanded by default")
	}
}

func TestRenderNested_ItemCounts(t *testing.T) {
	field := tags.Field{
		Name: "Items",
		Fields: []tags.Field{
			{Name: "Name", Type: "string"},
		},
		IsSlice:     true,
		ElementType: "struct",
	}

	value := []struct{ Name string }{
		{Name: "Item1"},
		{Name: "Item2"},
		{Name: "Item3"},
	}

	html := renderToString(renderInput(field, value))

	if !strings.Contains(html, "(3)") {
		t.Error("Should display item count in header")
	}
}

func TestRenderNested_Summaries(t *testing.T) {
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
	}{
		{Name: "User Service", Type: "HTTP"},
	}

	html := renderToString(renderInput(field, value))

	if !strings.Contains(html, "User Service") {
		t.Error("Summary should contain Name field")
	}
	if !strings.Contains(html, "slice__summary") {
		t.Error("Should have summary span")
	}
}

func TestRenderNested_DepthClasses(t *testing.T) {
	field := tags.Field{
		Name: "Level0",
		Fields: []tags.Field{
			{Name: "Name", Type: "string"},
		},
		IsSlice:     true,
		ElementType: "struct",
	}

	value := []struct{ Name string }{{Name: "Test"}}

	html := renderToString(renderInput(field, value))

	if !strings.Contains(html, "slice--depth-0") {
		t.Error("Should have depth-0 CSS class")
	}
}

func TestRenderNested_CompactItems(t *testing.T) {
	field := tags.Field{
		Name: "Services",
		Fields: []tags.Field{
			{Name: "Name", Type: "string"},
			{Name: "Endpoints", IsSlice: true, ElementType: "struct", Fields: []tags.Field{
				{Name: "Path", Type: "string"},
			}},
		},
		IsSlice:     true,
		ElementType: "struct",
	}

	value := []struct {
		Name      string
		Endpoints []struct{ Path string }
	}{
		{Name: "API", Endpoints: []struct{ Path string }{{Path: "/users"}}},
	}

	html := renderToString(renderInput(field, value))

	if !strings.Contains(html, "Endpoints") {
		t.Error("Should render nested Endpoints slice")
	}
}
