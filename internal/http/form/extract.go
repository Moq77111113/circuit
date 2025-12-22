package form

import (
	"reflect"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

// ExtractValues reads field values from a config struct.
// Recursively extracts all nested values with full dotted paths.
func ExtractValues(cfg any, s ast.Schema) map[string]any {
	values := make(map[string]any)
	rv := reflect.ValueOf(cfg).Elem()

	for _, node := range s.Nodes {
		fv := rv.FieldByName(node.Name)
		if !fv.IsValid() {
			continue
		}

		extractNodeValues(values, &node, fv, path.NewPath(node.Name))
	}

	return values
}

// extractNodeValues recursively extracts values for a node and its children.
func extractNodeValues(values map[string]any, node *ast.Node, fieldValue reflect.Value, currentPath path.Path) {
	val := fieldValue.Interface()
	if fieldValue.Kind() == reflect.Pointer && !fieldValue.IsNil() {
		val = fieldValue.Elem().Interface()
	}
	values[currentPath.String()] = val
	if node.Kind == ast.KindStruct && len(node.Children) > 0 {
		for _, child := range node.Children {
			childValue := fieldValue.FieldByName(child.Name)
			if !childValue.IsValid() {
				continue
			}
			childPath := currentPath.Child(child.Name)
			extractNodeValues(values, &child, childValue, childPath)
		}
	}
}
