package form

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/tags"
)

// ExtractValues reads field values from a config struct.
func ExtractValues(cfg any, s schema.Schema) map[string]any {
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
func Apply(cfg any, s schema.Schema, form url.Values) error {
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

		if field.IsSlice {
			if err := applySlice(fv, field, form); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
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

func applySlice(fv reflect.Value, field tags.Field, form url.Values) error {
	values := ParseIndexedField(form, field.Name)

	newSlice := reflect.MakeSlice(fv.Type(), len(values), len(values))

	applier, exists := appliers[field.ElementType]
	if !exists {
		return fmt.Errorf("no applier for type %s", field.ElementType)
	}

	for i, val := range values {
		if err := applier(newSlice.Index(i), val); err != nil {
			return fmt.Errorf("index %d: %w", i, err)
		}
	}

	fv.Set(newSlice)
	return nil
}
