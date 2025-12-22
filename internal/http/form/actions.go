package form

import (
	"fmt"
	"reflect"
	"strings"

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
// Handles dotted paths like "Database.Maintenance.AlertEmails".
func findNodeAndField(nodes []ast.Node, rootValue reflect.Value, path string) (*ast.Node, reflect.Value, error) {

	segments := strings.Split(path, ".")
	return findNodeAndFieldBySegments(nodes, rootValue, segments)
}

// findNodeAndFieldBySegments recursively traverses nodes and values following path segments.
func findNodeAndFieldBySegments(nodes []ast.Node, currentValue reflect.Value, segments []string) (*ast.Node, reflect.Value, error) {
	if len(segments) == 0 {
		return nil, reflect.Value{}, fmt.Errorf("empty path")
	}

	targetName := segments[0]
	for i := range nodes {
		if nodes[i].Name == targetName {
			fieldValue := currentValue.FieldByName(nodes[i].Name)
			if !fieldValue.IsValid() {
				return nil, reflect.Value{}, fmt.Errorf("field %s not found", targetName)
			}

			if len(segments) == 1 {
				return &nodes[i], fieldValue, nil
			}

			if nodes[i].Kind == ast.KindStruct && len(nodes[i].Children) > 0 {
				return findNodeAndFieldBySegments(nodes[i].Children, fieldValue, segments[1:])
			}

			return nil, reflect.Value{}, fmt.Errorf("cannot traverse into %s", targetName)
		}
	}

	return nil, reflect.Value{}, fmt.Errorf("field %s not found in schema", strings.Join(segments, "."))
}
