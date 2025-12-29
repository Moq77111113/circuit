package validation

import (
	"net/url"

	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// MergeVisitor merges form values with config values.
type MergeVisitor struct {
	form         url.Values
	configValues path.ValuesByPath
}

// VisitPrimitive merges form value for primitive fields.
func (v *MergeVisitor) VisitPrimitive(ctx *walk.VisitContext, n *node.Node) error {
	result := ctx.State.(path.ValuesByPath)
	fieldPath := ctx.Path.String()

	if formValue := v.form.Get(fieldPath); formValue != "" || v.form.Has(fieldPath) {
		result[fieldPath] = formValue
	} else if configValue, ok := v.configValues[fieldPath]; ok {
		result[fieldPath] = configValue
	}

	return nil
}

// VisitStruct handles struct nodes (recurses to children).
func (v *MergeVisitor) VisitStruct(ctx *walk.VisitContext, n *node.Node) error {
	return nil
}

// VisitSlice handles slice nodes (future implementation).
func (v *MergeVisitor) VisitSlice(ctx *walk.VisitContext, n *node.Node) error {
	return nil
}

// MergeFormValues merges form data with config values, preferring form values.
func MergeFormValues(nodes []node.Node, configValues path.ValuesByPath, form url.Values) path.ValuesByPath {
	result := make(path.ValuesByPath)

	visitor := &MergeVisitor{
		form:         form,
		configValues: configValues,
	}

	tree := &node.Tree{Nodes: nodes}
	walker := walk.NewWalker(visitor)
	_ = walker.Walk(tree, result)

	return result
}
