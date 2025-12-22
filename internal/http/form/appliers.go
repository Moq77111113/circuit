package form

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/moq77111113/circuit/internal/ast"
)

type valueApplier func(reflect.Value, string) error

var appliers = map[ast.ValueType]valueApplier{
	ast.ValueString: applyString,
	ast.ValueInt:    applyInt,
	ast.ValueBool:   applyBool,
	ast.ValueFloat:  applyFloat,
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

func applyFloat(fv reflect.Value, value string) error {
	if value == "" {
		fv.SetFloat(0)
		return nil
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid float: %w", err)
	}

	fv.SetFloat(val)
	return nil
}
