package validation

import (
	"fmt"
	"strings"

	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/path"
)

// validateRequired checks if a required field has a non-empty value.
func validateRequired(n *node.Node, value string, p path.Path) *ValidationError {
	if !n.UI.Required {
		return nil
	}
	if value == "" {
		return &ValidationError{
			Path:    p,
			Field:   n.Name,
			Message: fmt.Sprintf("%s is required", n.Name),
		}
	}
	return nil
}

// validateMinMax checks if numeric values respect min/max constraints.
func validateMinMax(n *node.Node, value string, p path.Path) *ValidationError {
	if value == "" {
		return nil
	}

	switch n.ValueType {
	case node.ValueInt:
		return validateIntMinMax(n, value, p)
	case node.ValueFloat:
		return validateFloatMinMax(n, value, p)
	default:
		return nil
	}
}

// validateSelectOptions checks if a select/radio value is in the allowed options.
func validateSelectOptions(n *node.Node, value string, p path.Path) *ValidationError {
	if n.UI.InputType != "select" && n.UI.InputType != "radio" {
		return nil
	}
	if value == "" {
		return nil
	}
	if len(n.UI.Options) == 0 {
		return nil
	}

	for _, opt := range n.UI.Options {
		if opt.Value == value {
			return nil
		}
	}

	var optionValues []string
	for _, opt := range n.UI.Options {
		optionValues = append(optionValues, opt.Value)
	}

	return &ValidationError{
		Path:    p,
		Field:   n.Name,
		Message: fmt.Sprintf("%s must be one of: %s", n.Name, strings.Join(optionValues, ", ")),
	}
}
