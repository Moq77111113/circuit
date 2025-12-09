package http

import (
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func parseIndexedField(form url.Values, fieldName string) []string {
	prefix := fieldName + "."
	indices := make(map[int]string)
	keys := []int{}

	for key, values := range form {
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		indexStr := strings.TrimPrefix(key, prefix)
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			continue
		}

		if len(values) > 0 {
			indices[index] = values[0]
			keys = append(keys, index)
		}
	}

	sort.Ints(keys)

	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, indices[key])
	}

	return result
}
