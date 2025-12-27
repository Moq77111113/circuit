package form

import (
	"net/url"
	"reflect"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

func ApplyNodes(cfg any, nodes []ast.Node, form url.Values) error {
	tree := &ast.Tree{Nodes: nodes}
	visitor := &FormVisitor{form: form}
	rv := reflect.ValueOf(cfg).Elem()

	for i := range nodes {
		node := &nodes[i]
		fieldValue := rv.FieldByName(node.Name)

		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			continue
		}

		ctx := &walk.VisitContext{
			Tree:   tree,
			State:  fieldValue,
			Path:   path.NewPath(node.Name),
			Depth:  0,
			Parent: nil,
			Index:  -1,
		}

		if err := visitor.dispatchNode(node, fieldValue, ctx); err != nil {
			return err
		}
	}

	return nil
}
