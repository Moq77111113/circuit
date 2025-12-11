package form

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"

	"github.com/moq77111113/circuit/internal/schema"
)

// applySliceNode applies form values to a slice field.
func applySliceNode(fv reflect.Value, node schema.Node, basePath schema.Path, form url.Values) error {
	if node.ElementKind == schema.KindStruct {
		return applyStructSlice(fv, node, basePath, form)
	}
	return applyPrimitiveSlice(fv, node, basePath, form)
}

// applyStructSlice applies form values to a slice of structs.
func applyStructSlice(fv reflect.Value, node schema.Node, basePath schema.Path, form url.Values) error {
	indices := extractSliceIndices(form, basePath)
	newSlice := reflect.MakeSlice(fv.Type(), len(indices), len(indices))

	for i, index := range indices {
		itemPath := basePath.Index(index)
		if err := applyNodes(newSlice.Index(i), node.Children, itemPath, form); err != nil {
			return fmt.Errorf("index %d: %w", index, err)
		}
	}

	fv.Set(newSlice)
	return nil
}

// applyPrimitiveSlice applies form values to a slice of primitives.
func applyPrimitiveSlice(fv reflect.Value, node schema.Node, basePath schema.Path, form url.Values) error {
	values := extractIndexedValues(form, basePath)
	newSlice := reflect.MakeSlice(fv.Type(), len(values), len(values))

	applier, exists := appliers[node.ValueType]
	if !exists {
		return fmt.Errorf("no applier for type %v", node.ValueType)
	}

	for i, val := range values {
		if err := applier(newSlice.Index(i), val); err != nil {
			return fmt.Errorf("index %d: %w", i, err)
		}
	}

	fv.Set(newSlice)
	return nil
}

// extractSliceIndices finds all unique slice indices in the form for a given path.
func extractSliceIndices(form url.Values, basePath schema.Path) []int {
	uniqueIndices := make(map[int]bool)

	for key := range form {
		path := schema.ParsePath(key)
		if path.HasPrefix(basePath) {
			if idx := path.IndexAfter(basePath); idx >= 0 {
				uniqueIndices[idx] = true
			}
		}
	}

	indices := make([]int, 0, len(uniqueIndices))
	for idx := range uniqueIndices {
		indices = append(indices, idx)
	}
	sort.Ints(indices)
	return indices
}

// extractIndexedValues extracts indexed values from the form for a primitive slice.
func extractIndexedValues(form url.Values, basePath schema.Path) []string {
	indices := extractSliceIndices(form, basePath)
	values := make([]string, len(indices))
	for i, idx := range indices {
		path := basePath.Index(idx)
		values[i] = form.Get(path.String())
	}
	return values
}
