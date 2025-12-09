package containers

import "fmt"

// DepthClass returns the CSS class for a given depth level.
func DepthClass(depth int) string {
	return fmt.Sprintf("slice--depth-%d", depth)
}

// IsCollapsed returns true if the item should be collapsed by default based on depth.
func IsCollapsed(depth int) bool {
	return depth >= 2
}
