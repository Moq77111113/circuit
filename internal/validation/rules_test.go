package validation

import (
	"testing"

	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/tags"
)

func TestValidateRequired(t *testing.T) {
	t.Run("required field empty returns error", func(t *testing.T) {
		n := &node.Node{
			Name:      "Host",
			ValueType: node.ValueString,
			UI: &node.UIMetadata{
				Required: true,
			},
		}
		p := path.NewPath("Host")

		err := validateRequired(n, "", p)
		if err == nil {
			t.Fatal("expected error for empty required field")
		}
		if err.Field != "Host" {
			t.Errorf("expected Field 'Host', got '%s'", err.Field)
		}
		if err.Message != "Host is required" {
			t.Errorf("expected Message 'Host is required', got '%s'", err.Message)
		}
	})

	t.Run("required field present returns nil", func(t *testing.T) {
		n := &node.Node{
			Name:      "Host",
			ValueType: node.ValueString,
			UI: &node.UIMetadata{
				Required: true,
			},
		}
		p := path.NewPath("Host")

		err := validateRequired(n, "localhost", p)
		if err != nil {
			t.Errorf("expected nil error for non-empty required field, got %v", err)
		}
	})

	t.Run("non-required field empty returns nil", func(t *testing.T) {
		n := &node.Node{
			Name:      "Host",
			ValueType: node.ValueString,
			UI: &node.UIMetadata{
				Required: false,
			},
		}
		p := path.NewPath("Host")

		err := validateRequired(n, "", p)
		if err != nil {
			t.Errorf("expected nil error for empty non-required field, got %v", err)
		}
	})
}

func TestValidateMinMax(t *testing.T) {
	t.Run("number below min returns error", func(t *testing.T) {
		n := &node.Node{
			Name:      "Port",
			ValueType: node.ValueInt,
			UI: &node.UIMetadata{
				Min: "1",
				Max: "65535",
			},
		}
		p := path.NewPath("Port")

		err := validateMinMax(n, "0", p)
		if err == nil {
			t.Fatal("expected error for value below min")
		}
		if err.Field != "Port" {
			t.Errorf("expected Field 'Port', got '%s'", err.Field)
		}
		if err.Message != "Port must be at least 1" {
			t.Errorf("expected Message 'Port must be at least 1', got '%s'", err.Message)
		}
	})

	t.Run("number above max returns error", func(t *testing.T) {
		n := &node.Node{
			Name:      "Port",
			ValueType: node.ValueInt,
			UI: &node.UIMetadata{
				Min: "1",
				Max: "65535",
			},
		}
		p := path.NewPath("Port")

		err := validateMinMax(n, "99999", p)
		if err == nil {
			t.Fatal("expected error for value above max")
		}
		if err.Field != "Port" {
			t.Errorf("expected Field 'Port', got '%s'", err.Field)
		}
		if err.Message != "Port must be at most 65535" {
			t.Errorf("expected Message 'Port must be at most 65535', got '%s'", err.Message)
		}
	})

	t.Run("number in range returns nil", func(t *testing.T) {
		n := &node.Node{
			Name:      "Port",
			ValueType: node.ValueInt,
			UI: &node.UIMetadata{
				Min: "1",
				Max: "65535",
			},
		}
		p := path.NewPath("Port")

		err := validateMinMax(n, "8080", p)
		if err != nil {
			t.Errorf("expected nil error for value in range, got %v", err)
		}
	})

	t.Run("no min/max constraints returns nil", func(t *testing.T) {
		n := &node.Node{
			Name:      "Port",
			ValueType: node.ValueInt,
			UI: &node.UIMetadata{
				Min: "",
				Max: "",
			},
		}
		p := path.NewPath("Port")

		err := validateMinMax(n, "99999", p)
		if err != nil {
			t.Errorf("expected nil error when no constraints, got %v", err)
		}
	})

	t.Run("float below min returns error", func(t *testing.T) {
		n := &node.Node{
			Name:      "Rate",
			ValueType: node.ValueFloat,
			UI: &node.UIMetadata{
				Min: "0.0",
				Max: "1.0",
			},
		}
		p := path.NewPath("Rate")

		err := validateMinMax(n, "-0.5", p)
		if err == nil {
			t.Fatal("expected error for float below min")
		}
		if err.Message != "Rate must be at least 0.0" {
			t.Errorf("expected Message 'Rate must be at least 0.0', got '%s'", err.Message)
		}
	})

	t.Run("float above max returns error", func(t *testing.T) {
		n := &node.Node{
			Name:      "Rate",
			ValueType: node.ValueFloat,
			UI: &node.UIMetadata{
				Min: "0.0",
				Max: "1.0",
			},
		}
		p := path.NewPath("Rate")

		err := validateMinMax(n, "1.5", p)
		if err == nil {
			t.Fatal("expected error for float above max")
		}
		if err.Message != "Rate must be at most 1.0" {
			t.Errorf("expected Message 'Rate must be at most 1.0', got '%s'", err.Message)
		}
	})

	t.Run("float in range returns nil", func(t *testing.T) {
		n := &node.Node{
			Name:      "Rate",
			ValueType: node.ValueFloat,
			UI: &node.UIMetadata{
				Min: "0.0",
				Max: "1.0",
			},
		}
		p := path.NewPath("Rate")

		err := validateMinMax(n, "0.5", p)
		if err != nil {
			t.Errorf("expected nil error for float in range, got %v", err)
		}
	})

	t.Run("non-numeric type returns nil", func(t *testing.T) {
		n := &node.Node{
			Name:      "Host",
			ValueType: node.ValueString,
			UI: &node.UIMetadata{
				Min: "1",
				Max: "100",
			},
		}
		p := path.NewPath("Host")

		err := validateMinMax(n, "localhost", p)
		if err != nil {
			t.Errorf("expected nil error for non-numeric type, got %v", err)
		}
	})
}

func TestValidateSelectOptions(t *testing.T) {
	t.Run("valid option returns nil", func(t *testing.T) {
		n := &node.Node{
			Name: "LogLevel",
			UI: &node.UIMetadata{
				InputType: tags.TypeSelect,
				Options: []tags.Option{
					{Value: "debug", Label: "Debug"},
					{Value: "info", Label: "Info"},
					{Value: "error", Label: "Error"},
				},
			},
		}
		p := path.NewPath("LogLevel")

		err := validateSelectOptions(n, "info", p)
		if err != nil {
			t.Errorf("expected nil error for valid option, got %v", err)
		}
	})

	t.Run("invalid option returns error", func(t *testing.T) {
		n := &node.Node{
			Name: "LogLevel",
			UI: &node.UIMetadata{
				InputType: tags.TypeSelect,
				Options: []tags.Option{
					{Value: "debug", Label: "Debug"},
					{Value: "info", Label: "Info"},
				},
			},
		}
		p := path.NewPath("LogLevel")

		err := validateSelectOptions(n, "invalid", p)
		if err == nil {
			t.Fatal("expected error for invalid option")
		}
		if err.Field != "LogLevel" {
			t.Errorf("expected Field 'LogLevel', got '%s'", err.Field)
		}
		if err.Message != "LogLevel must be one of: debug, info" {
			t.Errorf("expected option list message, got '%s'", err.Message)
		}
	})

	t.Run("non-select type returns nil", func(t *testing.T) {
		n := &node.Node{
			Name: "Host",
			UI: &node.UIMetadata{
				InputType: tags.TypeText,
			},
		}
		p := path.NewPath("Host")

		err := validateSelectOptions(n, "anything", p)
		if err != nil {
			t.Errorf("expected nil error for non-select type, got %v", err)
		}
	})

	t.Run("empty value with required returns nil for required validation", func(t *testing.T) {
		n := &node.Node{
			Name: "LogLevel",
			UI: &node.UIMetadata{
				InputType: tags.TypeSelect,
				Required:  true,
				Options: []tags.Option{
					{Value: "debug", Label: "Debug"},
				},
			},
		}
		p := path.NewPath("LogLevel")

		err := validateSelectOptions(n, "", p)
		if err != nil {
			t.Errorf("expected nil (defer to required validation), got %v", err)
		}
	})
}
