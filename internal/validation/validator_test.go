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

func TestValidate_StringLengthConstraints(t *testing.T) {
	schema := node.Schema{
		Name: "Config",
		Nodes: []node.Node{
			{
				Name:      "Username",
				Kind:      node.KindPrimitive,
				ValueType: node.ValueString,
				UI: &node.UIMetadata{
					MinLen: 3,
					MaxLen: 20,
				},
			},
		},
	}

	t.Run("value too short", func(t *testing.T) {
		form := url.Values{}
		form.Set("Username", "ab")

		result := Validate(schema, form)

		if result.Valid {
			t.Error("expected Valid to be false")
		}
		if len(result.Errors) != 1 {
			t.Fatalf("expected 1 error, got %d", len(result.Errors))
		}
		if result.Errors[0].Field != "Username" {
			t.Errorf("expected error for Username, got %s", result.Errors[0].Field)
		}
	})

	t.Run("value too long", func(t *testing.T) {
		form := url.Values{}
		form.Set("Username", "this_is_way_too_long_username")

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
		form.Set("Username", "validuser")

		result := Validate(schema, form)

		if !result.Valid {
			t.Errorf("expected Valid to be true, got errors: %v", result.Errors)
		}
		if len(result.Errors) != 0 {
			t.Errorf("expected 0 errors, got %d", len(result.Errors))
		}
	})
}

func TestValidate_PartialFormSubmission(t *testing.T) {
	// Schema has two fields: one at root level, one in nested struct
	// Simulates partial form submission when focused on a specific section
	schema := node.Schema{
		Name: "Config",
		Nodes: []node.Node{
			{
				Name:      "IntervalMS",
				Kind:      node.KindPrimitive,
				ValueType: node.ValueInt,
				UI: &node.UIMetadata{
					Min: "100",
					Max: "60000",
				},
			},
			{
				Name: "UDP",
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

	t.Run("only root fields submitted", func(t *testing.T) {
		// User is on main page with cards, only submits root-level fields
		form := url.Values{}
		form.Set("IntervalMS", "1000")
		// Note: UDP.Host is NOT in form (it's in a collapsed card)

		result := Validate(schema, form)

		if !result.Valid {
			t.Errorf("expected Valid to be true, got errors: %v", result.Errors)
		}
		if len(result.Errors) != 0 {
			t.Errorf("expected 0 errors (UDP.Host not in form), got %d", len(result.Errors))
		}
	})

	t.Run("nested fields submitted", func(t *testing.T) {
		// User is focused on UDP section, submits those fields
		form := url.Values{}
		form.Set("UDP.Host", "")
		form.Set("UDP.Port", "8080")
		// Note: IntervalMS is NOT in form

		result := Validate(schema, form)

		if result.Valid {
			t.Error("expected Valid to be false (UDP.Host is required and empty)")
		}
		if len(result.Errors) != 1 {
			t.Fatalf("expected 1 error (UDP.Host required), got %d", len(result.Errors))
		}
		if result.Errors[0].Path.String() != "UDP.Host" {
			t.Errorf("expected error path 'UDP.Host', got '%s'", result.Errors[0].Path.String())
		}
	})

	t.Run("full form submitted", func(t *testing.T) {
		// All fields present and valid
		form := url.Values{}
		form.Set("IntervalMS", "1000")
		form.Set("UDP.Host", "127.0.0.1")
		form.Set("UDP.Port", "8080")

		result := Validate(schema, form)

		if !result.Valid {
			t.Errorf("expected Valid to be true, got errors: %v", result.Errors)
		}
	})
}

func TestValidate_PatternConstraints(t *testing.T) {
	t.Run("email preset", func(t *testing.T) {
		schema := node.Schema{
			Name: "Config",
			Nodes: []node.Node{
				{
					Name:      "Email",
					Kind:      node.KindPrimitive,
					ValueType: node.ValueString,
					UI: &node.UIMetadata{
						Pattern: "email",
					},
				},
			},
		}

		t.Run("valid email", func(t *testing.T) {
			form := url.Values{}
			form.Set("Email", "user@example.com")

			result := Validate(schema, form)

			if !result.Valid {
				t.Errorf("expected Valid to be true, got errors: %v", result.Errors)
			}
		})

		t.Run("invalid email", func(t *testing.T) {
			form := url.Values{}
			form.Set("Email", "not-an-email")

			result := Validate(schema, form)

			if result.Valid {
				t.Error("expected Valid to be false")
			}
			if len(result.Errors) != 1 {
				t.Fatalf("expected 1 error, got %d", len(result.Errors))
			}
		})
	})

	t.Run("custom regex", func(t *testing.T) {
		schema := node.Schema{
			Name: "Config",
			Nodes: []node.Node{
				{
					Name:      "Code",
					Kind:      node.KindPrimitive,
					ValueType: node.ValueString,
					UI: &node.UIMetadata{
						Pattern: "^[A-Z]{2}[0-9]{4}$",
					},
				},
			},
		}

		t.Run("valid code", func(t *testing.T) {
			form := url.Values{}
			form.Set("Code", "AB1234")

			result := Validate(schema, form)

			if !result.Valid {
				t.Errorf("expected Valid to be true, got errors: %v", result.Errors)
			}
		})

		t.Run("invalid code", func(t *testing.T) {
			form := url.Values{}
			form.Set("Code", "invalid")

			result := Validate(schema, form)

			if result.Valid {
				t.Error("expected Valid to be false")
			}
			if len(result.Errors) != 1 {
				t.Fatalf("expected 1 error, got %d", len(result.Errors))
			}
		})
	})
}
