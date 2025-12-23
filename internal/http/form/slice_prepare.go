package form

import "reflect"

func prepareSlice(fieldValue reflect.Value, indices []int) reflect.Value {
	currentLen := fieldValue.Len()
	maxIdx := indices[len(indices)-1]
	targetLen := maxIdx + 1
	if currentLen > targetLen {
		targetLen = currentLen
	}

	newSlice := reflect.MakeSlice(fieldValue.Type(), targetLen, targetLen)
	for i := 0; i < currentLen; i++ {
		newSlice.Index(i).Set(fieldValue.Index(i))
	}
	return newSlice
}
