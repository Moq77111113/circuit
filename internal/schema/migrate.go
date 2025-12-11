package schema

import "github.com/moq77111113/circuit/internal/tags"

func FromTags(fields []tags.Field) []Node {
	nodes := make([]Node, 0, len(fields))
	for _, f := range fields {
		nodes = append(nodes, fromField(f))
	}
	return nodes
}

func fromField(f tags.Field) Node {
	n := Node{
		Name:      f.Name,
		Help:      f.Help,
		Required:  f.Required,
		Min:       f.Min,
		Max:       f.Max,
		Step:      f.Step,
		Options:   f.Options,
		InputType: f.InputType,
	}

	if f.IsSlice {
		n.Kind = KindSlice
		if len(f.Fields) > 0 {
			n.ElementKind = KindStruct
			n.Children = FromTags(f.Fields)
		} else {
			n.ElementKind = KindPrimitive
			n.ValueType = parseValueType(f.ElementType)
		}
	} else if f.InputType == tags.TypeSection {
		n.Kind = KindStruct
		n.Children = FromTags(f.Fields)
	} else {
		n.Kind = KindPrimitive
		n.ValueType = parseValueType(f.Type)
	}

	return n
}

func parseValueType(kindStr string) ValueType {
	switch kindStr {
	case "string":
		return ValueString
	case "int", "int64", "int32", "int16", "int8",
		"uint", "uint64", "uint32", "uint16", "uint8":
		return ValueInt
	case "bool":
		return ValueBool
	case "float64", "float32":
		return ValueFloat
	default:
		return ValueString
	}
}
