package tags

import "reflect"

// dereferenceType recursively dereferences pointer types until a non-pointer type is reached.
func dereferenceType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// elementType extracts the element type from a slice.
// Returns the element type and true if the input is a slice, otherwise returns the input type and false.
func elementType(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() == reflect.Slice {
		return t.Elem(), true
	}
	return t, false
}
