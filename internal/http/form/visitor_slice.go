package form

import (
	"reflect"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// applyStructSliceItems applies form values to struct items in a slice.
func (v *FormVisitor) applyStructSliceItems(ctx *walk.VisitContext, node *ast.Node, newSlice reflect.Value, indices []int) error {
	for i, idx := range indices {
		itemValue := newSlice.Index(i)
		itemPath := ctx.Path.Index(idx)

		// Apply each child field of the struct
		for _, child := range node.Children {
			childFieldValue := itemValue.FieldByName(child.Name)
			if !childFieldValue.IsValid() || !childFieldValue.CanSet() {
				continue
			}

			childPath := itemPath.Child(child.Name)

			// Create context for this child
			childState := &FormState{Current: childFieldValue}
			childCtx := &walk.VisitContext{
				Tree:   ctx.Tree,
				State:  childState,
				Path:   childPath,
				Depth:  ctx.Depth + 1,
				Parent: node,
				Index:  -1,
			}

			// Visit the child based on its kind
			switch child.Kind {
			case ast.KindPrimitive:
				if err := v.VisitPrimitive(childCtx, &child); err != nil {
					return err
				}
			case ast.KindStruct:
				// For nested structs in slices, we need to recurse
				childState.Current = childFieldValue
				for _, grandchild := range child.Children {
					grandchildFieldValue := childFieldValue.FieldByName(grandchild.Name)
					if !grandchildFieldValue.IsValid() || !grandchildFieldValue.CanSet() {
						continue
					}

					grandchildPath := childPath.Child(grandchild.Name)
					grandchildState := &FormState{Current: grandchildFieldValue}
					grandchildCtx := &walk.VisitContext{
						Tree:   ctx.Tree,
						State:  grandchildState,
						Path:   grandchildPath,
						Depth:  ctx.Depth + 2,
						Parent: &child,
						Index:  -1,
					}

					if grandchild.Kind == ast.KindPrimitive {
						if err := v.VisitPrimitive(grandchildCtx, &grandchild); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
