package tags

import (
	"errors"
	"reflect"
)

// Extract extracts fields from the struct tags of the given struct pointer.
func Extract(v any) ([]Field, error) {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Pointer {
		return nil, errors.New("extract: argument must be a pointer")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil, errors.New("extract: argument must be a pointer to struct")
	}

	rt := rv.Type()
	return extractFields(rt), nil
}

func extractFields(rt reflect.Type) []Field {
	var fields []Field

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("circuit")

		if tag == "-" {
			continue
		}

		f := Field{
			Name: field.Name,
			Type: field.Type.Kind().String(),
		}

		// Set default InputType based on Go type
		switch field.Type.Kind() {
		case reflect.Bool:
			f.InputType = TypeCheckbox
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			f.InputType = TypeNumber
		case reflect.String:
			f.InputType = TypeText
		}

		// Handle nested structs
		if field.Type.Kind() == reflect.Struct && field.Type.Name() != "Time" {
			f.Fields = extractFields(field.Type)
			f.InputType = TypeSection
		}

		parseTag(tag, &f)

		fields = append(fields, f)
	}

	return fields
}
