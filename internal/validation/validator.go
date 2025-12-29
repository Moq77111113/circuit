package validation

import (
	"net/url"

	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

// ValidationVisitor implements walk.Visitor to validate form data against schema.
type ValidationVisitor struct {
	form url.Values
}

// VisitPrimitive validates a primitive field.
func (v *ValidationVisitor) VisitPrimitive(ctx *walk.VisitContext, n *node.Node) error {
	result := ctx.State.(*ValidationResult)

	fieldPath := ctx.Path.String()
	value := v.form.Get(fieldPath)

	if err := validateRequired(n, value, ctx.Path); err != nil {
		result.Errors = append(result.Errors, *err)
		result.Valid = false
	}

	if err := validateMinMax(n, value, ctx.Path); err != nil {
		result.Errors = append(result.Errors, *err)
		result.Valid = false
	}

	if err := validateSelectOptions(n, value, ctx.Path); err != nil {
		result.Errors = append(result.Errors, *err)
		result.Valid = false
	}

	return nil
}

// VisitStruct validates a struct node (delegates to children).
func (v *ValidationVisitor) VisitStruct(ctx *walk.VisitContext, n *node.Node) error {
	return nil
}

// VisitSlice validates a slice node (future implementation).
func (v *ValidationVisitor) VisitSlice(ctx *walk.VisitContext, n *node.Node) error {
	return nil
}

// Validate validates form data against the schema and returns a ValidationResult.
func Validate(schema node.Schema, form url.Values) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	visitor := &ValidationVisitor{
		form: form,
	}

	tree := &node.Tree{Nodes: schema.Nodes}
	walker := walk.NewWalker(visitor)
	_ = walker.Walk(tree, result)

	return result
}
