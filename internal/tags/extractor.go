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

		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("circuit")

		if tag == "-" {
			continue
		}

		fieldType := dereferenceType(field.Type)
		elemType, isSlice := elementType(fieldType)

		if isSlice {
			fieldType = elemType
		}

		f := Field{
			Name:        field.Name,
			IsSlice:     isSlice,
			Type:        fieldType.Kind().String(),
			ElementType: fieldType.Kind().String(),
		}

		if fieldType.Kind() == reflect.Struct && fieldType.Name() != "Time" {
			f.Fields = extractFields(fieldType)
			if isSlice {
				f.Type = "slice"
			} else {
				f.Type = fieldType.Kind().String()
				f.InputType = TypeSection
			}
		} else {
			if isSlice {
				f.Type = "slice"
			}

			switch fieldType.Kind() {
			case reflect.Bool:
				f.InputType = TypeCheckbox
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				f.InputType = TypeNumber
			case reflect.String:
				f.InputType = TypeText
			}
		}

		parseTag(tag, &f)

		fields = append(fields, f)
	}

	return fields
}
