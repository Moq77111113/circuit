package yaml

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Encode encodes the given source structure into YAML format
func Encode(src any) ([]byte, error) {
	data, err := yaml.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("encode yaml: %w", err)
	}

	return data, nil
}
