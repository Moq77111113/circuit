package inputs

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
