package validation

import (
	"fmt"
	"strconv"

	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/path"
)

func validateIntMinMax(n *node.Node, value string, p path.Path) *ValidationError {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil
	}

	if n.UI.Min != "" {
		min, err := strconv.ParseInt(n.UI.Min, 10, 64)
		if err == nil && val < min {
			return &ValidationError{
				Path:    p,
				Field:   n.Name,
				Message: fmt.Sprintf("%s must be at least %s", n.Name, n.UI.Min),
			}
		}
	}

	if n.UI.Max != "" {
		max, err := strconv.ParseInt(n.UI.Max, 10, 64)
		if err == nil && val > max {
			return &ValidationError{
				Path:    p,
				Field:   n.Name,
				Message: fmt.Sprintf("%s must be at most %s", n.Name, n.UI.Max),
			}
		}
	}

	return nil
}

func validateFloatMinMax(n *node.Node, value string, p path.Path) *ValidationError {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}

	if n.UI.Min != "" {
		min, err := strconv.ParseFloat(n.UI.Min, 64)
		if err == nil && val < min {
			return &ValidationError{
				Path:    p,
				Field:   n.Name,
				Message: fmt.Sprintf("%s must be at least %s", n.Name, n.UI.Min),
			}
		}
	}

	if n.UI.Max != "" {
		max, err := strconv.ParseFloat(n.UI.Max, 64)
		if err == nil && val > max {
			return &ValidationError{
				Path:    p,
				Field:   n.Name,
				Message: fmt.Sprintf("%s must be at most %s", n.Name, n.UI.Max),
			}
		}
	}

	return nil
}
