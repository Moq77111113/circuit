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

			childCtx := &walk.VisitContext{
				Tree:   ctx.Tree,
				State:  childFieldValue,
				Path:   itemPath.Child(child.Name),
				Depth:  ctx.Depth + 1,
				Parent: node,
				Index:  -1,
			}

			if err := v.dispatchNode(&child, childFieldValue, childCtx); err != nil {
				return err
			}
		}
	}

	return nil
}
