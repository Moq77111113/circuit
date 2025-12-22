package path

import (
	"strconv"
	"strings"
)

type Path struct {
	segments []segment
}

type segment struct {
	name  string
	index int
}

func NewPath(name string) Path {
	return Path{segments: []segment{{name: name, index: -1}}}
}

func Root() Path {
	return Path{segments: []segment{}}
}

func (p Path) Child(name string) Path {
	p2 := Path{segments: make([]segment, len(p.segments)+1)}
	copy(p2.segments, p.segments)
	p2.segments[len(p.segments)] = segment{name: name, index: -1}
	return p2
}

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

func (p Path) IsRoot() bool {
	return len(p.segments) == 0
}

func (p Path) Segments() []string {
	if len(p.segments) == 0 {
		return nil
	}
	var parts []string
	for _, seg := range p.segments {
		parts = append(parts, seg.name)
		if seg.index >= 0 {
			parts = append(parts, strconv.Itoa(seg.index))
		}
	}
	return parts
}
