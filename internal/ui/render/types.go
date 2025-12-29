package render

import "github.com/moq77111113/circuit/internal/ast"

// valueTypeToString converts ValueType enum to string
func valueTypeToString(vt ast.ValueType) string {
	switch vt {
	case ast.ValueString:
		return "string"
	case ast.ValueInt:
		return "int"
	case ast.ValueBool:
		return "bool"
	case ast.ValueFloat:
		return "float"
	default:
		return "string"
	}
}
