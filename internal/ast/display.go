package ast

import "strings"

// DisplayName returns a human-readable display name for a node.
func DisplayName(n *Node) string {
	return SimplifyPath(n.Name)
}

// SimplifyPath extracts the last segment of a dot-separated path.
func SimplifyPath(pathStr string) string {
	if !strings.Contains(pathStr, ".") {
		return pathStr
	}
	parts := strings.Split(pathStr, ".")
	return parts[len(parts)-1]
}
