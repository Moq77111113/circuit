package node

import (
	"fmt"
	"reflect"

	"github.com/moq77111113/circuit/internal/tags"
)

// Extract builds a schema from a config struct.
// The argument must be a pointer to a struct.
func Extract(v any) (Schema, error) {
	fields, err := tags.Extract(v)
	if err != nil {
		return Schema{}, fmt.Errorf("schema extract: %w", err)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return Schema{}, fmt.Errorf("schema extract: expected pointer")
	}

	rv = rv.Elem()
	name := rv.Type().Name()

	return Schema{
		Name:  name,
		Nodes: FromTags(fields),
	}, nil
}
