package form

import (
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func ParseIndexedField(form url.Values, fieldName string) []string {
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

func ParseSliceIndices(form url.Values, fieldName string) []int {
	prefix := fieldName + "."
	uniqueIndices := make(map[int]bool)

	for key := range form {
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		rest := strings.TrimPrefix(key, prefix)
		parts := strings.SplitN(rest, ".", 2)
		indexStr := parts[0]

		index, err := strconv.Atoi(indexStr)
		if err != nil {
			continue
		}
		uniqueIndices[index] = true
	}

	keys := make([]int, 0, len(uniqueIndices))
	for k := range uniqueIndices {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}
