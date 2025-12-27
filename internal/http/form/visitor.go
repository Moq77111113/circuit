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

func (v *FormVisitor) dispatchNode(node *ast.Node, fieldValue reflect.Value, ctx *walk.VisitContext) error {
	ctx.State = fieldValue

	switch node.Kind {
	case ast.KindPrimitive:
		return v.VisitPrimitive(ctx, node)
	case ast.KindStruct:
		return v.VisitStruct(ctx, node)
	case ast.KindSlice:
		return v.VisitSlice(ctx, node)
	}
	return nil
}

func (v *FormVisitor) VisitPrimitive(ctx *walk.VisitContext, node *ast.Node) error {
	fieldValue := ctx.State.(reflect.Value)

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
	structValue := ctx.State.(reflect.Value)

	for i := range node.Children {
		child := &node.Children[i]
		childFieldValue := structValue.FieldByName(child.Name)

		if !childFieldValue.IsValid() || !childFieldValue.CanSet() {
			continue
		}

		childCtx := &walk.VisitContext{
			Tree:   ctx.Tree,
			State:  childFieldValue,
			Path:   ctx.Path.Child(child.Name),
			Depth:  ctx.Depth + 1,
			Parent: node,
			Index:  -1,
		}

		if err := v.dispatchNode(child, childFieldValue, childCtx); err != nil {
			return err
		}
	}

	return nil
}

func (v *FormVisitor) VisitSlice(ctx *walk.VisitContext, node *ast.Node) error {
	fieldValue := ctx.State.(reflect.Value)

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
