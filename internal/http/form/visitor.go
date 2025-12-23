package form

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

type FormVisitor struct {
	form url.Values
}

func (v *FormVisitor) VisitPrimitive(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*FormState)
	fieldValue := state.Current

	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		return nil
	}

	pathStr := ctx.Path.String()

	if !v.form.Has(pathStr) {
		return nil
	}

	formValue := v.form.Get(pathStr)

	applier, exists := appliers[node.ValueType]
	if !exists {
		return fmt.Errorf("no applier for type %v", node.ValueType)
	}

	return applier(fieldValue, formValue)
}

func (v *FormVisitor) VisitStruct(ctx *walk.VisitContext, node *ast.Node) error {
	return nil
}

func (v *FormVisitor) VisitSlice(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*FormState)
	fieldValue := state.Current

	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		return nil
	}

	indices := extractSliceIndices(v.form, ctx.Path)
	if len(indices) == 0 {
		if !hasAnyKeyWithPrefix(v.form, ctx.Path.String()+".") {
			return nil
		}
		fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 0))
		return nil
	}

	newSlice := prepareSlice(fieldValue, indices)

	if node.ElementKind == ast.KindPrimitive {
		if err := v.applyPrimitiveSliceItems(ctx, node, newSlice, indices); err != nil {
			return err
		}
	} else {
		if err := v.applyStructSliceItems(ctx, node, newSlice, indices); err != nil {
			return err
		}
	}

	fieldValue.Set(newSlice)
	return nil
}
