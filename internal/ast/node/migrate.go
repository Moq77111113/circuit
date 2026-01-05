package node

import (
	"github.com/moq77111113/circuit/internal/tags"
)

func FromTags(fields []tags.Field) []Node {
	nodes := make([]Node, 0, len(fields))
	for _, f := range fields {
		nodes = append(nodes, fromField(f))
	}
	return nodes
}

func fromField(f tags.Field) Node {
	n := Node{
		Name: f.Name,
		UI: &UIMetadata{
			InputType: f.InputType,
			Help:      f.Help,
			Required:  f.Required,
			ReadOnly:  f.ReadOnly,
			Min:       f.Min,
			Max:       f.Max,
			Step:      f.Step,
			Pattern:   f.Pattern,
			MinLen:    f.MinLen,
			MaxLen:    f.MaxLen,
			Options:   f.Options,
		},
	}

	if f.IsSlice {
		n.Kind = KindSlice
		if len(f.Fields) > 0 {
			n.ElementKind = KindStruct
			n.Children = FromTags(f.Fields)
		} else {
			n.ElementKind = KindPrimitive
			n.ValueType = ParseValueType(f.ElementType)
		}
	} else if f.InputType == tags.TypeSection {
		n.Kind = KindStruct
		n.Children = FromTags(f.Fields)
	} else {
		n.Kind = KindPrimitive
		n.ValueType = ParseValueType(f.Type)
	}

	return n
}
