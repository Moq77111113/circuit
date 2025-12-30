package toml

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

// Encode serializes src into TOML format.
func Encode(src any) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(src)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
