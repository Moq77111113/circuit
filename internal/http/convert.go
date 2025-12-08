package http

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/tags"
)

type valueApplier func(reflect.Value, string) error

var appliers = map[string]valueApplier{
	"string": applyString,
	"int":    applyInt,
	"bool":   applyBool,
}

// extractValues reads field values from a config struct.
func extractValues(cfg any, s schema.Schema) map[string]any {
	values := make(map[string]any)
	rv := reflect.ValueOf(cfg).Elem()

	for _, field := range s.Fields {
		fv := rv.FieldByName(field.Name)
		if !fv.IsValid() {
			continue
		}

		values[field.Name] = fv.Interface()
	}

	return values
}

// applyForm updates a config struct from form data.
func applyForm(cfg any, s schema.Schema, form url.Values) error {
	rv := reflect.ValueOf(cfg).Elem()
	return applyValues(rv, s.Fields, form)
}

func applyValues(rv reflect.Value, fields []tags.Field, form url.Values) error {
	for _, field := range fields {
		fv := rv.FieldByName(field.Name)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}

		if field.InputType == tags.TypeSection {
			if err := applyValues(fv, field.Fields, form); err != nil {
				return err
			}
			continue
		}

		applier, exists := appliers[field.Type]
		if !exists {
			continue
		}

		if err := applier(fv, form.Get(field.Name)); err != nil {
			return fmt.Errorf("%s: %w", field.Name, err)
		}
	}

	return nil
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
