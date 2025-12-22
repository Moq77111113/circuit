package node

// NodeKind represents the kind of schema node
type NodeKind uint8

const (
	KindPrimitive NodeKind = iota // string, int, bool, float
	KindStruct                    // nested object
	KindSlice                     // []T
)

// ValueType represents the primitive value type
type ValueType uint8

const (
	ValueString ValueType = iota
	ValueInt
	ValueBool
	ValueFloat
)
