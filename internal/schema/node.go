package schema

import "github.com/moq77111113/circuit/internal/tags"

// NodeKind represents the kind of schema node
type NodeKind uint8

const (
	KindPrimitive NodeKind = iota // string, int, bool, float
	KindStruct                    // nested object
	KindSlice                     // []T
)

type ValueType uint8

const (
	ValueString ValueType = iota
	ValueInt
	ValueBool
	ValueFloat
)

type Node struct {
	Name        string
	Kind        NodeKind
	ValueType   ValueType
	InputType   tags.InputType
	Children    []Node
	ElementKind NodeKind

	Help        string
	Placeholder string
	Required    bool
	Min         string
	Max         string
	Step        string
	Options     []tags.Option
}
