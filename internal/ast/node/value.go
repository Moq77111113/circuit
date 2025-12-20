package node

// ParseValueType converts a Go type string to ValueType.
func ParseValueType(kindStr string) ValueType {
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
