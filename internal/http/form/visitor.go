package form

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// FormVisitor implements walk.Visitor for form data parsing.
type FormVisitor struct {
	form url.Values
}

// VisitPrimitive applies a form value to a primitive field.
func (v *FormVisitor) VisitPrimitive(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*FormState)
	fieldValue := state.Current

	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		return nil
	}

	formValue := v.form.Get(ctx.Path.String())

	applier, exists := appliers[node.ValueType]
	if !exists {
		return fmt.Errorf("no applier for type %v", node.ValueType)
	}

	return applier(fieldValue, formValue)
}

// VisitStruct applies form values to a struct's fields.
// The Walker handles recursion into children automatically.
func (v *FormVisitor) VisitStruct(ctx *walk.VisitContext, node *ast.Node) error {
	// Nothing to do - Walker will visit all children
	return nil
}

// VisitSlice applies form values to a slice.
func (v *FormVisitor) VisitSlice(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*FormState)
	fieldValue := state.Current

	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		return nil
	}

	// Extract slice indices from form keys
	indices := extractSliceIndices(v.form, ctx.Path)
	if len(indices) == 0 {
		// No items in form - set empty slice
		fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 0))
		return nil
	}

	newSlice := reflect.MakeSlice(fieldValue.Type(), len(indices), len(indices))

	if node.ElementKind == ast.KindPrimitive {
		// Primitive slice: apply values directly
		for i, idx := range indices {
			itemPath := ctx.Path.Index(idx)
			formValue := v.form.Get(itemPath.String())

			applier := appliers[node.ValueType]
			if applier == nil {
				return fmt.Errorf("no applier for primitive slice type %v", node.ValueType)
			}

			if err := applier(newSlice.Index(i), formValue); err != nil {
				return fmt.Errorf("slice item %d: %w", idx, err)
			}
		}
	} else {
		// Struct slice: apply children for each item
		if err := v.applyStructSliceItems(ctx, node, newSlice, indices); err != nil {
			return err
		}
	}

	fieldValue.Set(newSlice)
	return nil
}
