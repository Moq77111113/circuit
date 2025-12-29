package render

import (
	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ui/styles"
	"github.com/moq77111113/circuit/internal/validation"
)

// RenderContext bundles all data and configuration needed for rendering.
type RenderContext struct {
	// Schema + data
	Schema *ast.Schema
	Values ast.ValuesByPath
	Focus  path.Path

	// Render behavior
	CollapseDepthThreshold int
	ShowCardsAtDepth0      bool
	MaxDepth               int
	ReadOnly               bool
	Errors                 *validation.ValidationResult
}

// NewRenderContext creates a RenderContext with sensible defaults.
func NewRenderContext(s *ast.Schema, values ast.ValuesByPath) *RenderContext {
	return &RenderContext{
		Schema:                 s,
		Values:                 values,
		Focus:                  path.Root(),
		CollapseDepthThreshold: 2,
		ShowCardsAtDepth0:      true,
		MaxDepth:               4,
		ReadOnly:               false,
	}
}

// ShouldCollapse returns true if items at the given depth should be collapsed.
func (rc *RenderContext) ShouldCollapse(depth int) bool {
	return depth >= rc.CollapseDepthThreshold
}

// ClampDepth returns the depth clamped to the range [0, MaxDepth].
func (rc *RenderContext) ClampDepth(depth int) int {
	if depth < 0 {
		return 0
	}
	if depth > rc.MaxDepth {
		return rc.MaxDepth
	}
	return depth
}

// DepthClass returns the class for the given depth (clamped).
func (rc *RenderContext) DepthClass(depth int) string {
	return styles.DepthClass(rc.ClampDepth(depth))
}
