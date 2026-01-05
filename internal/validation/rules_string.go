package validation

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/path"
)

// validateMinLen checks if string length >= MinLen.
func validateMinLen(n *node.Node, value string, p path.Path) *ValidationError {
	if n.UI.MinLen == 0 {
		return nil
	}

	length := utf8.RuneCountInString(value)
	if length < n.UI.MinLen {
		return &ValidationError{
			Path:    p,
			Field:   n.Name,
			Message: fmt.Sprintf("must be at least %d characters (got %d)", n.UI.MinLen, length),
		}
	}
	return nil
}

// validateMaxLen checks if string length <= MaxLen.
func validateMaxLen(n *node.Node, value string, p path.Path) *ValidationError {
	if n.UI.MaxLen == 0 {
		return nil
	}

	length := utf8.RuneCountInString(value)
	if length > n.UI.MaxLen {
		return &ValidationError{
			Path:    p,
			Field:   n.Name,
			Message: fmt.Sprintf("must be at most %d characters (got %d)", n.UI.MaxLen, length),
		}
	}
	return nil
}

// validatePattern checks regex pattern or preset.
func validatePattern(n *node.Node, value string, p path.Path) *ValidationError {
	if n.UI.Pattern == "" {
		return nil
	}

	pattern := resolvePattern(n.UI.Pattern)

	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return &ValidationError{
			Path:    p,
			Field:   n.Name,
			Message: fmt.Sprintf("invalid pattern configuration: %v", err),
		}
	}

	if !matched {
		return &ValidationError{
			Path:    p,
			Field:   n.Name,
			Message: formatPatternError(n.UI.Pattern),
		}
	}
	return nil
}

// resolvePattern converts preset names to regex.
func resolvePattern(pattern string) string {
	if preset, ok := GetPreset(pattern); ok {
		return preset
	}
	return pattern
}

// formatPatternError returns user-friendly message.
func formatPatternError(pattern string) string {
	switch pattern {
	case "email":
		return "must be a valid email address"
	case "url":
		return "must be a valid URL"
	case "phone":
		return "must be a valid phone number"
	default:
		return fmt.Sprintf("must match pattern: %s", pattern)
	}
}
