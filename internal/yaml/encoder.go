package yaml

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Encode marshals a struct into YAML format.
// src can be a pointer or a value.
func Encode(src any) ([]byte, error) {
	data, err := yaml.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("encode yaml: %w", err)
	}

	return data, nil
}
