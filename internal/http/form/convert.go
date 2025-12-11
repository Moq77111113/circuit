package form

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/moq77111113/circuit/internal/schema"
)

// ExtractValues reads field values from a config struct
func ExtractValues(cfg any, s schema.Schema) map[string]any {
	values := make(map[string]any)
	rv := reflect.ValueOf(cfg).Elem()

	for _, node := range s.Nodes {
		fv := rv.FieldByName(node.Name)
		if !fv.IsValid() {
			continue
		}

		val := fv.Interface()
		if fv.Kind() == reflect.Pointer && !fv.IsNil() {
			val = fv.Elem().Interface()
		}
		values[node.Name] = val
	}

	return values
}

// Apply updates a config struct from form data
func Apply(cfg any, s schema.Schema, form url.Values) error {
	return ApplyNodes(cfg, s.Nodes, form)
}

// ApplyNodes updates a config struct from form data using the new Node system.
func ApplyNodes(cfg any, nodes []schema.Node, form url.Values) error {
	rv := reflect.ValueOf(cfg).Elem()
	return applyNodes(rv, nodes, schema.Path{}, form)
}

// applyNodes recursively applies form values to struct fields using Node and Path
func applyNodes(rv reflect.Value, nodes []schema.Node, basePath schema.Path, form url.Values) error {
	for _, node := range nodes {
		path := basePath.Child(node.Name)
		fv := rv.FieldByName(node.Name)

		if !fv.IsValid() || !fv.CanSet() {
			continue
		}

		switch node.Kind {
		case schema.KindPrimitive:
			if err := applyPrimitive(fv, node, path.String(), form); err != nil {
				return fmt.Errorf("%s: %w", node.Name, err)
			}

		case schema.KindStruct:
			if err := applyNodes(fv, node.Children, path, form); err != nil {
				return err
			}

		case schema.KindSlice:
			if err := applySliceNode(fv, node, path, form); err != nil {
				return fmt.Errorf("%s: %w", node.Name, err)
			}
		}
	}

	return nil
}

// applyPrimitive applies a form value to a primitive field.
func applyPrimitive(fv reflect.Value, node schema.Node, fieldName string, form url.Values) error {
	value := form.Get(fieldName)

	applier, exists := appliers[node.ValueType]
	if !exists {
		return fmt.Errorf("no applier for type %v", node.ValueType)
	}

	return applier(fv, value)
}
