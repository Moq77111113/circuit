package form

import (
	"fmt"
	"reflect"
	"strconv"
)

type valueApplier func(reflect.Value, string) error

var appliers = map[string]valueApplier{
	"string": applyString,
	"int":    applyInt,
	"bool":   applyBool,
}

func applyString(fv reflect.Value, value string) error {
	fv.SetString(value)
	return nil
}

func applyInt(fv reflect.Value, value string) error {
	if value == "" {
		fv.SetInt(0)
		return nil
	}

	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid int: %w", err)
	}

	fv.SetInt(val)
	return nil
}

func applyBool(fv reflect.Value, value string) error {
	fv.SetBool(value == "on")
	return nil
}
