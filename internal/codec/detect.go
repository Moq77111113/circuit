package codec

import (
	"fmt"
	"path/filepath"
)

var registry = make(map[Extension]Codec)

// Register associates a file extension with a codec implementation.
// This is typically called from codec package init() functions.
func Register(ext Extension, codec Codec) {
	registry[ext] = codec
}

// Detect returns the appropriate codec for the given file path
// based on its extension. Returns an error if the extension is
// not supported.
func Detect(path string) (Codec, error) {
	ext := filepath.Ext(path)
	codec, ok := registry[Extension(ext)]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}
	return codec, nil
}
