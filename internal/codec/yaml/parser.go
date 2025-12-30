package yaml

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

// Parse unmarshals YAML data into the destination struct.
// dst must be a pointer to a struct.
func Parse(data []byte, dst any) error {
	if reflect.TypeOf(dst).Kind() != reflect.Pointer {
		return fmt.Errorf("parse: destination must be a pointer")
	}

	if err := yaml.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	return nil
}
