package codec

import (
	"fmt"
	"path/filepath"
)

var registry = make(map[Extension]Codec)

// Register registers a Codec implementation for the given file extension
func Register(ext Extension, codec Codec) {
	registry[ext] = codec
}

// Detect returns the Codec implementation for the given file path
func Detect(path string) (Codec, error) {
	ext := filepath.Ext(path)
	codec, ok := registry[Extension(ext)]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}
	return codec, nil
}
