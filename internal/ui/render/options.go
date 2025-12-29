package render

import (
	"github.com/moq77111113/circuit/internal/ui/styles"
	"github.com/moq77111113/circuit/internal/validation"
)

// Options configures rendering behavior for the UI.
type Options struct {
	// CollapseDepthThreshold sets the depth at which collapsibles are collapsed by default.
	// Example: depth 2 means slices at depth 2 and below start collapsed.
	// Default: 2
	CollapseDepthThreshold int

	// ShowCardsAtDepth0 controls whether depth-0 structs render as clickable cards.
	// When true, top-level structs show preview cards with navigation.
	// Default: true
	ShowCardsAtDepth0 bool

	// MaxDepth limits the maximum nesting depth rendered.
	// Depths beyond this value are clamped to MaxDepth.
	// Default: 4
	MaxDepth int

	// Errors holds validation errors to display inline with fields.
	// If nil, no validation errors are shown.
	Errors *validation.ValidationResult

	// ReadOnly controls whether inputs are rendered as disabled.
	// When true, all inputs are disabled (readonly mode).
	// Default: false
	ReadOnly bool
}

// DefaultOptions returns sensible defaults for rendering.
func DefaultOptions() Options {
	return Options{
		CollapseDepthThreshold: 2,
		ShowCardsAtDepth0:      true,
		MaxDepth:               4,
	}
}

// ShouldCollapse returns true if items at the given depth should be collapsed by default.
func (o Options) ShouldCollapse(depth int) bool {
	return depth >= o.CollapseDepthThreshold
}

// ClampDepth returns the depth clamped to the range [0, MaxDepth].
func (o Options) ClampDepth(depth int) int {
	if depth < 0 {
		return 0
	}
	if depth > o.MaxDepth {
		return o.MaxDepth
	}
	return depth
}

// DepthClass returns the BEM CSS class for the given depth (clamped).
func (o Options) DepthClass(depth int) string {
	return styles.DepthClass(o.ClampDepth(depth))
}
