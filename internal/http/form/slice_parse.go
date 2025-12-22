package form

import (
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/moq77111113/circuit/internal/ast/path"
)

// extractSliceIndices parses form keys to find all slice indices for a given path.
// For example, if form has "Services.0.Name" and "Services.2.Name", returns [0, 2].
func extractSliceIndices(form url.Values, path path.Path) []int {
	pathPrefix := path.String() + "."
	indexMap := make(map[int]bool)

	for key := range form {
		if !strings.HasPrefix(key, pathPrefix) {
			continue
		}

		// Extract index from key like "Services.0.Name" → 0
		rest := strings.TrimPrefix(key, pathPrefix)
		parts := strings.SplitN(rest, ".", 2)
		if len(parts) == 0 {
			continue
		}

		idx, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		indexMap[idx] = true
	}

	// Convert map to sorted slice
	indices := make([]int, 0, len(indexMap))
	for idx := range indexMap {
		indices = append(indices, idx)
	}
	sort.Ints(indices)

	return indices
}
