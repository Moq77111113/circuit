package codec

// Codec defines the interface for configuration format parsing and encoding.
type Codec interface {
	Parse(data []byte, dst any) error
	Encode(src any) ([]byte, error)
}

// Extension represents a file extension for configuration formats.
type Extension string

const (
	ExtYAML Extension = ".yaml"
	ExtYML  Extension = ".yml"
	ExtTOML Extension = ".toml"
	ExtJSON Extension = ".json"
)
