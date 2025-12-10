package containers

import (
	"testing"
)

func TestDepthClass(t *testing.T) {
	tests := []struct {
		depth    int
		expected string
	}{
		{0, "slice--depth-0"},
		{1, "slice--depth-1"},
		{2, "slice--depth-2"},
		{10, "slice--depth-10"},
	}

	for _, tt := range tests {
		got := DepthClass(tt.depth)
		if got != tt.expected {
			t.Errorf("DepthClass(%d) = %q, want %q", tt.depth, got, tt.expected)
		}
	}
}

func TestIsCollapsed(t *testing.T) {
	tests := []struct {
		depth    int
		expected bool
	}{
		{0, false},
		{1, false},
		{2, true},
		{3, true},
	}

	for _, tt := range tests {
		got := IsCollapsed(tt.depth)
		if got != tt.expected {
			t.Errorf("IsCollapsed(%d) = %v, want %v", tt.depth, got, tt.expected)
		}
	}
}
