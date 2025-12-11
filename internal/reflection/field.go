package reflection

import "reflect"

// FieldByName extracts a field value from a struct by name using reflection.
func FieldByName(structValue any, fieldName string) any {
	if structValue == nil {
		return nil
	}
	v := reflect.ValueOf(structValue)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		return nil
	}
	return f.Interface()
}
