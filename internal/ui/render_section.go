package ui

import (
	"reflect"

	g "maragu.dev/gomponents"

	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/inputs"
)

func renderSection(field tags.Field, value any) g.Node {
	var children []g.Node

	for _, subField := range field.Fields {
		var subValue any
		if value != nil {
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Pointer {
				v = v.Elem()
			}
			if v.Kind() == reflect.Struct {
				f := v.FieldByName(subField.Name)
				if f.IsValid() {
					subValue = f.Interface()
				}
			}
		}
		children = append(children, renderField(subField, map[string]any{subField.Name: subValue}))
	}

	return inputs.Section(field.Name, children)
}
