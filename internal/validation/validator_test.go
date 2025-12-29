package validation

import (
	"net/url"
	"testing"

	"github.com/moq77111113/circuit/internal/ast/node"
)

func TestValidate_SingleRequiredFieldViolation(t *testing.T) {
	schema := node.Schema{
		Name: "Config",
		Nodes: []node.Node{
			{
				Name:      "Host",
				Kind:      node.KindPrimitive,
				ValueType: node.ValueString,
				UI: &node.UIMetadata{
					Required: true,
				},
			},
		},
	}

	form := url.Values{}
	form.Set("Host", "")

	result := Validate(schema, form)

	if result.Valid {
		t.Error("expected Valid to be false")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Field != "Host" {
		t.Errorf("expected error for 'Host', got '%s'", result.Errors[0].Field)
	}
	if result.Errors[0].Message != "Host is required" {
		t.Errorf("expected message 'Host is required', got '%s'", result.Errors[0].Message)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	schema := node.Schema{
		Name: "Config",
		Nodes: []node.Node{
			{
				Name:      "Host",
				Kind:      node.KindPrimitive,
				ValueType: node.ValueString,
				UI: &node.UIMetadata{
					Required: true,
				},
			},
			{
				Name:      "Port",
				Kind:      node.KindPrimitive,
				ValueType: node.ValueInt,
				UI: &node.UIMetadata{
					Required: true,
					Min:      "1",
					Max:      "65535",
				},
			},
		},
	}

	form := url.Values{}
	form.Set("Host", "")
	form.Set("Port", "99999")

	result := Validate(schema, form)

	if result.Valid {
		t.Error("expected Valid to be false")
	}
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(result.Errors))
	}
}

func TestValidate_NestedStructValidation(t *testing.T) {
	schema := node.Schema{
		Name: "Config",
		Nodes: []node.Node{
			{
				Name: "Server",
				Kind: node.KindStruct,
				Children: []node.Node{
					{
						Name:      "Host",
						Kind:      node.KindPrimitive,
						ValueType: node.ValueString,
						UI: &node.UIMetadata{
							Required: true,
						},
					},
					{
						Name:      "Port",
						Kind:      node.KindPrimitive,
						ValueType: node.ValueInt,
						UI: &node.UIMetadata{
							Min: "1",
							Max: "65535",
						},
					},
				},
			},
		},
	}

	form := url.Values{}
	form.Set("Server.Host", "")
	form.Set("Server.Port", "0")

	result := Validate(schema, form)

	if result.Valid {
		t.Error("expected Valid to be false")
	}
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(result.Errors))
	}

	hostErr := result.Errors[0]
	if hostErr.Field != "Host" {
		t.Errorf("expected first error Field 'Host', got '%s'", hostErr.Field)
	}
	if hostErr.Path.String() != "Server.Host" {
		t.Errorf("expected first error Path 'Server.Host', got '%s'", hostErr.Path.String())
	}

	portErr := result.Errors[1]
	if portErr.Field != "Port" {
		t.Errorf("expected second error Field 'Port', got '%s'", portErr.Field)
	}
	if portErr.Path.String() != "Server.Port" {
		t.Errorf("expected second error Path 'Server.Port', got '%s'", portErr.Path.String())
	}
}

func TestValidate_ValidForm(t *testing.T) {
	schema := node.Schema{
		Name: "Config",
		Nodes: []node.Node{
			{
				Name:      "Host",
				Kind:      node.KindPrimitive,
				ValueType: node.ValueString,
				UI: &node.UIMetadata{
					Required: true,
				},
			},
			{
				Name:      "Port",
				Kind:      node.KindPrimitive,
				ValueType: node.ValueInt,
				UI: &node.UIMetadata{
					Min: "1",
					Max: "65535",
				},
			},
		},
	}

	form := url.Values{}
	form.Set("Host", "localhost")
	form.Set("Port", "8080")

	result := Validate(schema, form)

	if !result.Valid {
		t.Error("expected Valid to be true")
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d", len(result.Errors))
	}
}

func TestValidate_OptionalFieldEmpty(t *testing.T) {
	schema := node.Schema{
		Name: "Config",
		Nodes: []node.Node{
			{
				Name:      "Host",
				Kind:      node.KindPrimitive,
				ValueType: node.ValueString,
				UI: &node.UIMetadata{
					Required: false,
				},
			},
		},
	}

	form := url.Values{}
	form.Set("Host", "")

	result := Validate(schema, form)

	if !result.Valid {
		t.Error("expected Valid to be true for empty optional field")
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d", len(result.Errors))
	}
}

func TestValidate_MinMaxConstraints(t *testing.T) {
	schema := node.Schema{
		Name: "Config",
		Nodes: []node.Node{
			{
				Name:      "Port",
				Kind:      node.KindPrimitive,
				ValueType: node.ValueInt,
				UI: &node.UIMetadata{
					Min: "1024",
					Max: "49151",
				},
			},
		},
	}

	t.Run("value below min", func(t *testing.T) {
		form := url.Values{}
		form.Set("Port", "80")

		result := Validate(schema, form)

		if result.Valid {
			t.Error("expected Valid to be false")
		}
		if len(result.Errors) != 1 {
			t.Fatalf("expected 1 error, got %d", len(result.Errors))
		}
	})

	t.Run("value above max", func(t *testing.T) {
		form := url.Values{}
		form.Set("Port", "65535")

		result := Validate(schema, form)

		if result.Valid {
			t.Error("expected Valid to be false")
		}
		if len(result.Errors) != 1 {
			t.Fatalf("expected 1 error, got %d", len(result.Errors))
		}
	})

	t.Run("value in range", func(t *testing.T) {
		form := url.Values{}
		form.Set("Port", "8080")

		result := Validate(schema, form)

		if !result.Valid {
			t.Error("expected Valid to be true")
		}
		if len(result.Errors) != 0 {
			t.Errorf("expected 0 errors, got %d", len(result.Errors))
		}
	})
}
