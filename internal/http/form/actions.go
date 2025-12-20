package form

import (
	"fmt"
	"reflect"

	"github.com/moq77111113/circuit/internal/ast"
)

// AddSliceItemNode adds a new item to a slice field in the config.
func AddSliceItemNode(cfg any, nodes []ast.Node, fieldPath string) error {
	rv := reflect.ValueOf(cfg).Elem()

	// Find the node and field
	node, fieldValue, err := findNodeAndField(nodes, rv, fieldPath)
	if err != nil {
		return err
	}

	if node.Kind != ast.KindSlice {
		return fmt.Errorf("%s is not a slice", fieldPath)
	}

	// Create a new zero-value item
	elemType := fieldValue.Type().Elem()
	newItem := reflect.New(elemType).Elem()

	// Append to slice
	newSlice := reflect.Append(fieldValue, newItem)
	fieldValue.Set(newSlice)

	return nil
}

// RemoveSliceItemNode removes an item from a slice field in the config.
func RemoveSliceItemNode(cfg any, nodes []ast.Node, fieldPath string, index int) error {
	rv := reflect.ValueOf(cfg).Elem()

	// Find the node and field
	node, fieldValue, err := findNodeAndField(nodes, rv, fieldPath)
	if err != nil {
		return err
	}

	if node.Kind != ast.KindSlice {
		return fmt.Errorf("%s is not a slice", fieldPath)
	}

	if index < 0 || index >= fieldValue.Len() {
		return fmt.Errorf("index %d out of range for slice of length %d", index, fieldValue.Len())
	}

	// Remove item by creating a new slice without it
	newSlice := reflect.MakeSlice(fieldValue.Type(), 0, fieldValue.Len()-1)
	for i := 0; i < fieldValue.Len(); i++ {
		if i != index {
			newSlice = reflect.Append(newSlice, fieldValue.Index(i))
		}
	}

	fieldValue.Set(newSlice)
	return nil
}

// findNodeAndField finds a node and its corresponding field value by path.
func findNodeAndField(nodes []ast.Node, rootValue reflect.Value, path string) (*ast.Node, reflect.Value, error) {
	// Parse path (e.g., "Services" or "Server.Database")
	currentValue := rootValue
	currentNodes := nodes

	// For now, assume simple paths (just field name, no dots)
	// This matches the current usage in the handler
	for _, node := range currentNodes {
		if node.Name == path {
			fieldValue := currentValue.FieldByName(node.Name)
			if !fieldValue.IsValid() {
				return nil, reflect.Value{}, fmt.Errorf("field %s not found", path)
			}
			return &node, fieldValue, nil
		}

		// Check nested structs
		if node.Kind == ast.KindStruct {
			// Recursively search in children
			childValue := currentValue.FieldByName(node.Name)
			if childValue.IsValid() {
				foundNode, foundValue, err := findNodeAndField(node.Children, childValue, path)
				if err == nil {
					return foundNode, foundValue, nil
				}
			}
		}
	}

	return nil, reflect.Value{}, fmt.Errorf("field %s not found in schema", path)
}
