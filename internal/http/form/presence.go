package form

import (
	"net/url"
	"strings"
)

func hasAnyKeyWithPrefix(form url.Values, prefix string) bool {
	for key := range form {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}
