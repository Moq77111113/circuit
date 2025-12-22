package path

import (
	"strconv"
	"strings"
)

// ParsePath extracts a path from form field name like "Services.0.Endpoints.2.Name".
func ParsePath(fieldName string) Path {
	if fieldName == "" {
		return Path{}
	}

	parts := strings.Split(fieldName, ".")
	var segments []segment

	for _, part := range parts {
		if idx, err := strconv.Atoi(part); err == nil {
			if len(segments) > 0 {
				segments[len(segments)-1].index = idx
			}
		} else {
			segments = append(segments, segment{name: part, index: -1})
		}
	}

	return Path{segments: segments}
}

// HasPrefix checks if path starts with prefix.
func (p Path) HasPrefix(prefix Path) bool {
	if len(prefix.segments) > len(p.segments) {
		return false
	}
	for i, prefixSeg := range prefix.segments {
		pathSeg := p.segments[i]
		if pathSeg.name != prefixSeg.name {
			return false
		}
		if prefixSeg.index >= 0 && pathSeg.index != prefixSeg.index {
			return false
		}
	}
	return true
}

// IndexAfter returns the index of the last prefix segment in the path.
func (p Path) IndexAfter(prefix Path) int {
	if !p.HasPrefix(prefix) {
		return -1
	}
	if len(prefix.segments) == 0 {
		return -1
	}
	lastPrefixIdx := len(prefix.segments) - 1
	return p.segments[lastPrefixIdx].index
}
