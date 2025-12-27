package validation

import (
	"testing"

	"github.com/moq77111113/circuit/internal/ast/path"
)

func TestValidationError(t *testing.T) {
	p := path.NewPath("Server").Child("Port")
	err := ValidationError{
		Path:    p,
		Field:   "Port",
		Message: "Port is required",
	}

	if err.Field != "Port" {
		t.Errorf("expected Field 'Port', got '%s'", err.Field)
	}
	if err.Message != "Port is required" {
		t.Errorf("expected Message 'Port is required', got '%s'", err.Message)
	}
	if err.Path.String() != "Server.Port" {
		t.Errorf("expected Path 'Server.Port', got '%s'", err.Path.String())
	}
}

func TestValidationResult_Has(t *testing.T) {
	p1 := path.NewPath("Host")
	p2 := path.NewPath("Port")
	p3 := path.NewPath("Timeout")

	result := &ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{Path: p1, Field: "Host", Message: "Host is required"},
			{Path: p2, Field: "Port", Message: "Port must be between 1 and 65535"},
		},
	}

	if !result.Has(p1) {
		t.Error("expected Has(p1) to return true")
	}
	if !result.Has(p2) {
		t.Error("expected Has(p2) to return true")
	}
	if result.Has(p3) {
		t.Error("expected Has(p3) to return false")
	}
}

func TestValidationResult_Get(t *testing.T) {
	p1 := path.NewPath("Host")
	p2 := path.NewPath("Port")
	p3 := path.NewPath("Timeout")

	result := &ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{Path: p1, Field: "Host", Message: "Host is required"},
			{Path: p2, Field: "Port", Message: "Port must be between 1 and 65535"},
		},
	}

	if msg := result.Get(p1); msg != "Host is required" {
		t.Errorf("expected message 'Host is required', got '%s'", msg)
	}
	if msg := result.Get(p2); msg != "Port must be between 1 and 65535" {
		t.Errorf("expected message 'Port must be between 1 and 65535', got '%s'", msg)
	}
	if msg := result.Get(p3); msg != "" {
		t.Errorf("expected empty message for p3, got '%s'", msg)
	}
}

func TestValidationResult_FirstError(t *testing.T) {
	t.Run("with errors", func(t *testing.T) {
		p1 := path.NewPath("Host")
		p2 := path.NewPath("Port")

		result := &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{Path: p1, Field: "Host", Message: "Host is required"},
				{Path: p2, Field: "Port", Message: "Port must be between 1 and 65535"},
			},
		}

		first := result.FirstError()
		if first == nil {
			t.Fatal("expected non-nil first error")
		}
		if first.Field != "Host" {
			t.Errorf("expected first error Field 'Host', got '%s'", first.Field)
		}
		if first.Message != "Host is required" {
			t.Errorf("expected first error Message 'Host is required', got '%s'", first.Message)
		}
	})

	t.Run("without errors", func(t *testing.T) {
		result := &ValidationResult{
			Valid:  true,
			Errors: []ValidationError{},
		}

		first := result.FirstError()
		if first != nil {
			t.Errorf("expected nil first error, got %v", first)
		}
	})
}

func TestValidationResult_Valid(t *testing.T) {
	t.Run("valid when no errors", func(t *testing.T) {
		result := &ValidationResult{
			Valid:  true,
			Errors: []ValidationError{},
		}

		if !result.Valid {
			t.Error("expected Valid to be true when no errors")
		}
	})

	t.Run("invalid when has errors", func(t *testing.T) {
		p := path.NewPath("Host")
		result := &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{Path: p, Field: "Host", Message: "Host is required"},
			},
		}

		if result.Valid {
			t.Error("expected Valid to be false when has errors")
		}
	})
}
