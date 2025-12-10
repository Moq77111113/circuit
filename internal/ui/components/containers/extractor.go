package containers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/moq77111113/circuit/internal/tags"
)

type Summary struct {
	Fields []Field
}

type Field struct {
	Name  string
	Value string
}

func Extract(field tags.Field, value any, maxFields int) Summary {
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
	for i := 0; i < len(field.Fields) && len(fields) < maxFields; i++ {
		f := field.Fields[i]
		fv := v.FieldByName(f.Name)
		if !fv.IsValid() {
			continue
		}

		valueStr := extractFieldValue(f, fv)
		if valueStr == "" {
			continue
		}
		fields = append(fields, Field{Name: f.Name, Value: valueStr})
	}

	return Summary{Fields: fields}
}

func extractFieldValue(f tags.Field, fv reflect.Value) string {
	switch f.Type {
	case "string":
		return fv.String()
	case "bool":
		if fv.Bool() {
			return "true"
		}
	case "int":
		if i := fv.Int(); i != 0 {
			return fmt.Sprintf("%d", i)
		}
	case "slice":
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
