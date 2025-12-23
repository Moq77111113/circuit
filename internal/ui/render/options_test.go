package render

import (
	"testing"

	"github.com/moq77111113/circuit/internal/ui/styles"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if opts.CollapseDepthThreshold != 2 {
		t.Errorf("expected CollapseDepthThreshold=2, got %d", opts.CollapseDepthThreshold)
	}
	if opts.ShowCardsAtDepth0 != true {
		t.Errorf("expected ShowCardsAtDepth0=true, got %v", opts.ShowCardsAtDepth0)
	}
	if opts.MaxDepth != 4 {
		t.Errorf("expected MaxDepth=4, got %d", opts.MaxDepth)
	}
}

func TestShouldCollapse(t *testing.T) {
	tests := []struct {
		name      string
		threshold int
		depth     int
		want      bool
	}{
		{"depth 0 with threshold 2", 2, 0, false},
		{"depth 1 with threshold 2", 2, 1, false},
		{"depth 2 with threshold 2", 2, 2, true},
		{"depth 3 with threshold 2", 2, 3, true},
		{"depth 0 with threshold 0", 0, 0, true},
		{"depth 1 with threshold 3", 3, 1, false},
		{"depth 3 with threshold 3", 3, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{CollapseDepthThreshold: tt.threshold}
			got := opts.ShouldCollapse(tt.depth)
			if got != tt.want {
				t.Errorf("ShouldCollapse(%d) = %v, want %v", tt.depth, got, tt.want)
			}
		})
	}
}

func TestClampDepth(t *testing.T) {
	tests := []struct {
		name     string
		maxDepth int
		depth    int
		want     int
	}{
		{"negative clamped to 0", 4, -1, 0},
		{"depth 0", 4, 0, 0},
		{"depth 2", 4, 2, 2},
		{"depth 4 at max", 4, 4, 4},
		{"depth 5 clamped to max 4", 4, 5, 4},
		{"depth 10 clamped to max 4", 4, 10, 4},
		{"depth 3 with max 2", 2, 3, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{MaxDepth: tt.maxDepth}
			got := opts.ClampDepth(tt.depth)
			if got != tt.want {
				t.Errorf("ClampDepth(%d) = %d, want %d", tt.depth, got, tt.want)
			}
		})
	}
}

func TestDepthClass(t *testing.T) {
	opts := Options{MaxDepth: 4}

	tests := []struct {
		depth int
		want  string
	}{
		{0, styles.CollapsibleDepth0},
		{1, styles.CollapsibleDepth1},
		{2, styles.CollapsibleDepth2},
		{3, styles.CollapsibleDepth3},
		{4, styles.CollapsibleDepth4},
		{5, styles.CollapsibleDepth4},
		{10, styles.CollapsibleDepth4},
	}

	for _, tt := range tests {
		t.Run("depth "+string(rune(tt.depth+'0')), func(t *testing.T) {
			got := opts.DepthClass(tt.depth)
			if got != tt.want {
				t.Errorf("DepthClass(%d) = %q, want %q", tt.depth, got, tt.want)
			}
		})
	}
}
