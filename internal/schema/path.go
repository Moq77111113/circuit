package schema

import (
	"strconv"
	"strings"
)

// Path represents a hierarchical field path like "Services.0.Endpoints.2.Name".
type Path struct {
	segments []segment
}

type segment struct {
	name  string
	index int
}

// NewPath creates a root path with a single name.
func NewPath(name string) Path {
	return Path{segments: []segment{{name: name, index: -1}}}
}

// Child appends a struct field to the path.
func (p Path) Child(name string) Path {
	p2 := Path{segments: make([]segment, len(p.segments)+1)}
	copy(p2.segments, p.segments)
	p2.segments[len(p.segments)] = segment{name: name, index: -1}
	return p2
}

// Index appends a slice index to the last segment.
func (p Path) Index(idx int) Path {
	if len(p.segments) == 0 {
		return p
	}
	p2 := Path{segments: make([]segment, len(p.segments))}
	copy(p2.segments, p.segments)
	p2.segments[len(p2.segments)-1].index = idx
	return p2
}

func (p Path) String() string {
	if len(p.segments) == 0 {
		return ""
	}

	var parts []string
	for _, seg := range p.segments {
		parts = append(parts, seg.name)
		if seg.index >= 0 {
			parts = append(parts, strconv.Itoa(seg.index))
		}
	}
	return strings.Join(parts, ".")
}

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
