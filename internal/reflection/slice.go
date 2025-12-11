package reflection

import "reflect"

// SliceValues extracts all values from a slice using reflection.
func SliceValues(value any) []any {
	if value == nil {
		return nil
	}
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil
	}
	items := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		items[i] = v.Index(i).Interface()
	}
	return items
}
