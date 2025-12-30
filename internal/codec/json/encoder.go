package json

import "encoding/json"

// Encode serializes src into indented JSON format.
func Encode(src any) ([]byte, error) {
	return json.MarshalIndent(src, "", "  ")
}
