package form

import (
	"reflect"
)

// FormState tracks the current reflect.Value during form parsing traversal.
type FormState struct {
	Current reflect.Value // Current field value being processed
}

// NewFormState creates a new form parsing state with the root value.
func NewFormState(rootValue reflect.Value) *FormState {
	return &FormState{
		Current: rootValue,
	}
}
