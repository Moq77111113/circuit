package form

import (
	"fmt"
	"reflect"

	"github.com/moq77111113/circuit/internal/schema"
)

// AddSliceItemNode adds a zero-value item to a slice field using Node system.
func AddSliceItemNode(cfg any, nodes []schema.Node, fieldPath string) error {
	path := schema.ParsePath(fieldPath)
	node, rv := findNodeAndValue(nodes, path, reflect.ValueOf(cfg).Elem())

	if node == nil || node.Kind != schema.KindSlice {
		return fmt.Errorf("field %s not found or not a slice", fieldPath)
	}

	if !rv.IsValid() || !rv.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldPath)
	}

	elemType := rv.Type().Elem()
	zero := reflect.Zero(elemType)
	newSlice := reflect.Append(rv, zero)
	rv.Set(newSlice)

	return nil
}

// RemoveSliceItemNode removes an item at index from a slice field using Node system.
func RemoveSliceItemNode(cfg any, nodes []schema.Node, fieldPath string, index int) error {
	path := schema.ParsePath(fieldPath)
	node, rv := findNodeAndValue(nodes, path, reflect.ValueOf(cfg).Elem())

	if node == nil || node.Kind != schema.KindSlice {
		return fmt.Errorf("field %s not found or not a slice", fieldPath)
	}

	if !rv.IsValid() || !rv.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldPath)
	}

	if index < 0 || index >= rv.Len() {
		return fmt.Errorf("index %d out of range for field %s", index, fieldPath)
	}

	newSlice := reflect.AppendSlice(
		rv.Slice(0, index),
		rv.Slice(index+1, rv.Len()),
	)
	rv.Set(newSlice)

	return nil
}
