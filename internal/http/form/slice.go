package form

import (
	"fmt"
	"reflect"

	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/tags"
)

func AddSliceItem(cfg any, s schema.Schema, fieldName string) error {
	field := findField(s.Fields, fieldName)
	if field == nil || !field.IsSlice {
		return fmt.Errorf("field %s not found or not a slice", fieldName)
	}

	rv := reflect.ValueOf(cfg).Elem()
	fv := rv.FieldByName(field.Name)
	if !fv.IsValid() || !fv.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldName)
	}

	elemType := fv.Type().Elem()
	zero := reflect.Zero(elemType)
	newSlice := reflect.Append(fv, zero)
	fv.Set(newSlice)

	return nil
}

func RemoveSliceItem(cfg any, s schema.Schema, fieldName string, index int) error {
	field := findField(s.Fields, fieldName)
	if field == nil || !field.IsSlice {
		return fmt.Errorf("field %s not found or not a slice", fieldName)
	}

	rv := reflect.ValueOf(cfg).Elem()
	fv := rv.FieldByName(field.Name)
	if !fv.IsValid() || !fv.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldName)
	}

	if index < 0 || index >= fv.Len() {
		return fmt.Errorf("index %d out of range for field %s", index, fieldName)
	}

	newSlice := reflect.AppendSlice(
		fv.Slice(0, index),
		fv.Slice(index+1, fv.Len()),
	)
	fv.Set(newSlice)

	return nil
}

func findField(fields []tags.Field, name string) *tags.Field {
	for i := range fields {
		if fields[i].Name == name {
			return &fields[i]
		}
		if fields[i].InputType == tags.TypeSection {
			if found := findField(fields[i].Fields, name); found != nil {
				return found
			}
		}
	}
	return nil
}
