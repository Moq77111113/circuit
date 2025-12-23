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

func ExtractFromMap(children []ast.Node, itemMap map[string]any, maxFields int) Summary {
	if itemMap == nil || maxFields <= 0 {
		return Summary{}
	}

	fields := []Field{}
	for i := 0; i < len(children) && len(fields) < maxFields; i++ {
		child := children[i]
		if child.Kind != ast.KindPrimitive {
			continue
		}

		value := itemMap[child.Name]
		valueStr := formatValue(child.ValueType, value)

		fields = append(fields, Field{
			Name:  child.Name,
			Value: valueStr,
		})
	}

	return Summary{Fields: fields}
}

func formatValue(vt ast.ValueType, value any) string {
	if value == nil {
		return ""
	}

	switch vt {
	case ast.ValueString:
		return fmt.Sprintf("%v", value)
	case ast.ValueInt:
		return fmt.Sprintf("%d", value)
	case ast.ValueBool:
		return fmt.Sprintf("%t", value)
	case ast.ValueFloat:
		return fmt.Sprintf("%.2f", value)
	default:
		return ""
	}
}
