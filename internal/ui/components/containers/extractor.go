package containers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/moq77111113/circuit/internal/ast"
)

type Summary struct {
	Fields []Field
}

type Field struct {
	Name  string
	Value string
}

func Extract(node ast.Node, value any, maxFields int) Summary {
	if value == nil || maxFields <= 0 {
		return Summary{}
	}

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return Summary{}
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return Summary{}
	}

	var fields []Field
	for i := 0; i < len(node.Children) && len(fields) < maxFields; i++ {
		child := node.Children[i]
		fv := v.FieldByName(child.Name)
		if !fv.IsValid() {
			continue
		}

		valueStr := extractFieldValue(child, fv)
		if valueStr == "" {
			continue
		}
		fields = append(fields, Field{Name: child.Name, Value: valueStr})
	}

	return Summary{Fields: fields}
}

func extractFieldValue(node ast.Node, fv reflect.Value) string {
	switch node.Kind {
	case ast.KindPrimitive:
		switch node.ValueType {
		case ast.ValueString:
			return fv.String()
		case ast.ValueBool:
			if fv.Bool() {
				return "true"
			}
		case ast.ValueInt:
			if i := fv.Int(); i != 0 {
				return fmt.Sprintf("%d", i)
			}
		case ast.ValueFloat:
			if f := fv.Float(); f != 0 {
				return fmt.Sprintf("%.2f", f)
			}
		}
	case ast.KindSlice:
		if n := fv.Len(); n > 0 {
			return fmt.Sprintf("%d", n)
		}
	}
	return ""
}

func Format(s Summary) string {
	if len(s.Fields) == 0 {
		return ""
	}

	var parts []string
	for _, f := range s.Fields {
		parts = append(parts, fmt.Sprintf("%s: %s", f.Name, f.Value))
	}
	return strings.Join(parts, " • ")
}
