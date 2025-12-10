package containers

import "reflect"

func reflectSliceValues(value any) []any {
	if value == nil {
		return []any{}
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice {
		return []any{}
	}

	items := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		items[i] = rv.Index(i).Interface()
	}

	return items
}

func extractStructField(v any, fieldName string) any {
	if v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil
	}

	fv := rv.FieldByName(fieldName)
	if !fv.IsValid() {
		return nil
	}

	return fv.Interface()
}
