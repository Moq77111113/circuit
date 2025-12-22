package render

import (
	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

// FilterByFocus returns nodes to display based on the current focus path.
func FilterByFocus(nodes []ast.Node, focus path.Path) []ast.Node {
	if focus.IsRoot() {
		return stripStructChildren(nodes)
	}

	segments := focus.Segments()
	return findNodeByPath(nodes, segments)
}

func stripStructChildren(nodes []ast.Node) []ast.Node {
	result := make([]ast.Node, len(nodes))
	for i, node := range nodes {
		result[i] = node
		if node.Kind == ast.KindStruct {
			result[i].Children = nil
		}
	}
	return result
}

// findNodeByPath traverses the node tree following the path segments.
func findNodeByPath(nodes []ast.Node, segments []string) []ast.Node {
	if len(segments) == 0 {
		return nodes
	}

	currentSeg := segments[0]
	remaining := segments[1:]

	for _, node := range nodes {
		if node.Name == currentSeg {
			if len(remaining) == 0 {
				if node.Kind == ast.KindPrimitive {
					return []ast.Node{node}
				}
				if node.Kind == ast.KindSlice {
					return []ast.Node{node}
				}
				return stripStructChildren(node.Children)
			}
			if node.Kind == ast.KindStruct {
				return findNodeByPath(node.Children, remaining)
			}

			if node.Kind == ast.KindSlice && len(remaining) > 0 {

				return findNodeByPath(node.Children, remaining[1:])
			}
		}
	}


	return nodes
}

func IsRootFocus(focus path.Path) bool {
	return focus.IsRoot()
}

// ShouldRenderAsCard determines if a struct should be rendered as a clickable card
func ShouldRenderAsCard(node ast.Node, focus path.Path, nodePath path.Path) bool {

	if node.Kind != ast.KindStruct {
		return false
	}

	if focus.HasPrefix(nodePath) {
		return false
	}

	return true
}
