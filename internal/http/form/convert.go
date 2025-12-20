package form

import (
	"net/url"
	"reflect"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// ExtractValues reads field values from a config struct.
func ExtractValues(cfg any, s ast.Schema) map[string]any {
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

// Apply updates a config struct from form data using the Visitor pattern.
func Apply(cfg any, s ast.Schema, form url.Values) error {
	return ApplyNodes(cfg, s.Nodes, form)
}

// ApplyNodes updates a config struct from form data using FormVisitor.
func ApplyNodes(cfg any, nodes []ast.Node, form url.Values) error {
	tree := &ast.Tree{Nodes: nodes}
	visitor := &FormVisitor{form: form}

	rv := reflect.ValueOf(cfg).Elem()
	state := NewFormState(rv)

	// Walk the tree, visiting each node
	for i := range nodes {
		node := &nodes[i]
		fieldValue := rv.FieldByName(node.Name)

		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			continue
		}

		// Update state with current field
		state.Current = fieldValue

		// Create context for this node
		ctx := &walk.VisitContext{
			Tree:   tree,
			State:  state,
			Path:   path.NewPath(node.Name),
			Depth:  0,
			Parent: nil,
			Index:  -1,
		}

		// Visit based on kind
		var err error
		switch node.Kind {
		case ast.KindPrimitive:
			err = visitor.VisitPrimitive(ctx, node)
		case ast.KindStruct:
			// For structs, we need to recurse into children
			err = visitStructChildren(visitor, ctx, node, fieldValue)
		case ast.KindSlice:
			err = visitor.VisitSlice(ctx, node)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// visitStructChildren recursively visits all children of a struct node.
func visitStructChildren(visitor *FormVisitor, parentCtx *walk.VisitContext, node *ast.Node, structValue reflect.Value) error {
	for i := range node.Children {
		child := &node.Children[i]
		childFieldValue := structValue.FieldByName(child.Name)

		if !childFieldValue.IsValid() || !childFieldValue.CanSet() {
			continue
		}

		childState := &FormState{Current: childFieldValue}
		childCtx := &walk.VisitContext{
			Tree:   parentCtx.Tree,
			State:  childState,
			Path:   parentCtx.Path.Child(child.Name),
			Depth:  parentCtx.Depth + 1,
			Parent: node,
			Index:  -1,
		}

		var err error
		switch child.Kind {
		case ast.KindPrimitive:
			err = visitor.VisitPrimitive(childCtx, child)
		case ast.KindStruct:
			err = visitStructChildren(visitor, childCtx, child, childFieldValue)
		case ast.KindSlice:
			err = visitor.VisitSlice(childCtx, child)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
