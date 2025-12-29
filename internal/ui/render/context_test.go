package render

import (
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/validation"
)

func TestNewRenderContext(t *testing.T) {
	schema := &ast.Schema{Name: "TestConfig"}
	values := ast.ValuesByPath{"host": "localhost"}

	rc := NewRenderContext(schema, values)

	if rc.Schema != schema {
		t.Error("Schema not set correctly")
	}
	if rc.Values == nil || rc.Values["host"] != "localhost" {
		t.Error("Values not set correctly")
	}
	if !rc.Focus.IsRoot() {
		t.Error("Focus should default to root")
	}
	if rc.CollapseDepthThreshold != 2 {
		t.Errorf("Expected CollapseDepthThreshold=2, got %d", rc.CollapseDepthThreshold)
	}
	if !rc.ShowCardsAtDepth0 {
		t.Error("ShowCardsAtDepth0 should default to true")
	}
	if rc.MaxDepth != 4 {
		t.Errorf("Expected MaxDepth=4, got %d", rc.MaxDepth)
	}
	if rc.ReadOnly {
		t.Error("ReadOnly should default to false")
	}
	if rc.Errors != nil {
		t.Error("Errors should default to nil")
	}
}

func TestRenderContextShouldCollapse(t *testing.T) {
	tests := []struct {
		threshold int
		depth     int
		want      bool
	}{
		{threshold: 2, depth: 0, want: false},
		{threshold: 2, depth: 1, want: false},
		{threshold: 2, depth: 2, want: true},
		{threshold: 2, depth: 3, want: true},
		{threshold: 0, depth: 0, want: true},
		{threshold: 5, depth: 4, want: false},
	}

	for _, tt := range tests {
		rc := &RenderContext{CollapseDepthThreshold: tt.threshold}
		got := rc.ShouldCollapse(tt.depth)
		if got != tt.want {
			t.Errorf("ShouldCollapse(depth=%d, threshold=%d) = %v, want %v",
				tt.depth, tt.threshold, got, tt.want)
		}
	}
}

func TestRenderContextClampDepth(t *testing.T) {
	rc := &RenderContext{MaxDepth: 4}

	tests := []struct {
		input int
		want  int
	}{
		{input: -1, want: 0},
		{input: 0, want: 0},
		{input: 2, want: 2},
		{input: 4, want: 4},
		{input: 5, want: 4},
		{input: 10, want: 4},
	}

	for _, tt := range tests {
		got := rc.ClampDepth(tt.input)
		if got != tt.want {
			t.Errorf("ClampDepth(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestRenderContextDepthClass(t *testing.T) {
	rc := &RenderContext{MaxDepth: 4}

	tests := []struct {
		depth int
		want  string
	}{
		{depth: 0, want: "collapsible--depth-0"},
		{depth: 2, want: "collapsible--depth-2"},
		{depth: 4, want: "collapsible--depth-4"},
		{depth: 10, want: "collapsible--depth-4"}, // Clamped
	}

	for _, tt := range tests {
		got := rc.DepthClass(tt.depth)
		if got != tt.want {
			t.Errorf("DepthClass(%d) = %q, want %q", tt.depth, got, tt.want)
		}
	}
}

func TestRenderContextWithErrors(t *testing.T) {
	schema := &ast.Schema{Name: "Test"}
	values := ast.ValuesByPath{}

	rc := NewRenderContext(schema, values)

	// Set errors
	errors := &validation.ValidationResult{
		Errors: []validation.ValidationError{
			{Path: path.NewPath("field1"), Message: "Required field"},
		},
	}
	rc.Errors = errors

	if rc.Errors == nil {
		t.Error("Errors should be set")
	}

	errorMsg := rc.Errors.Get(path.NewPath("field1"))
	if errorMsg != "Required field" {
		t.Errorf("Expected error 'Required field', got %q", errorMsg)
	}
}

func TestRenderContextWithReadOnly(t *testing.T) {
	schema := &ast.Schema{Name: "Test"}
	values := ast.ValuesByPath{}

	rc := NewRenderContext(schema, values)
	rc.ReadOnly = true

	if !rc.ReadOnly {
		t.Error("ReadOnly should be true")
	}
}
