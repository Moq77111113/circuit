package yaml

import "github.com/moq77111113/circuit/internal/codec"

// Codec implements the codec.Codec interface for YAML format.
type Codec struct{}

func (c Codec) Parse(data []byte, dst any) error {
	return Parse(data, dst)
}

func (c Codec) Encode(src any) ([]byte, error) {
	return Encode(src)
}

func init() {
	codec.Register(codec.ExtYAML, Codec{})
	codec.Register(codec.ExtYML, Codec{})
}
