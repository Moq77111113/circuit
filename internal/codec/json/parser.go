package json

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Parse decodes JSON data into dst.
func Parse(data []byte, dst any) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}
	return json.Unmarshal(data, dst)
}
