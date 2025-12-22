package path

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

// String converts the path to dot-notation (e.g., "Services.0.Endpoints.2.Name").
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

// FieldPath returns the path without indices (e.g., "Services.Endpoints").
func (p Path) FieldPath() string {
	if len(p.segments) == 0 {
		return ""
	}

	var names []string
	for _, seg := range p.segments {
		names = append(names, seg.name)
	}
	return strings.Join(names, ".")
}
