package form

import (
	"fmt"
	"reflect"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

func (v *FormVisitor) applyPrimitiveSliceItems(ctx *walk.VisitContext, node *ast.Node, newSlice reflect.Value, indices []int) error {
	for _, idx := range indices {
		itemPath := ctx.Path.Index(idx)
		formValue := v.form.Get(itemPath.String())

		applier := appliers[node.ValueType]
		if applier == nil {
			return fmt.Errorf("no applier for primitive slice type %v", node.ValueType)
		}

		if err := applier(newSlice.Index(idx), formValue); err != nil {
			return fmt.Errorf("slice item %d: %w", idx, err)
		}
	}
	return nil
}

func (v *FormVisitor) applyStructSliceItems(ctx *walk.VisitContext, node *ast.Node, newSlice reflect.Value, indices []int) error {
	for _, idx := range indices {
		itemValue := newSlice.Index(idx)
		itemPath := ctx.Path.Index(idx)

		for _, child := range node.Children {
			childFieldValue := itemValue.FieldByName(child.Name)
			if !childFieldValue.IsValid() || !childFieldValue.CanSet() {
				continue
			}

			childPath := itemPath.Child(child.Name)

			childState := &FormState{Current: childFieldValue}
			childCtx := &walk.VisitContext{
				Tree:   ctx.Tree,
				State:  childState,
				Path:   childPath,
				Depth:  ctx.Depth + 1,
				Parent: node,
				Index:  -1,
			}

			switch child.Kind {
			case ast.KindPrimitive:
				if err := v.VisitPrimitive(childCtx, &child); err != nil {
					return err
				}
			case ast.KindStruct:
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
