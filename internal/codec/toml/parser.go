package toml

import (
	"fmt"
	"reflect"

	"github.com/BurntSushi/toml"
)

// Parse decodes TOML data into dst.
func Parse(data []byte, dst any) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}
	_, err := toml.Decode(string(data), dst)
	return err
}
