package validation

import (
	"testing"

	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/path"
)

func TestValidateMinLen(t *testing.T) {
	tests := []struct {
		name    string
		minLen  int
		value   string
		wantErr bool
	}{
		{"no constraint", 0, "hi", false},
		{"valid", 3, "hello", false},
		{"exact match", 5, "hello", false},
		{"too short", 10, "hello", true},
		{"unicode aware - café", 3, "café", false},
		{"unicode aware - 日本語", 3, "日本語", false},
		{"empty with minlen 0", 0, "", false},
		{"empty with minlen 1", 1, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &node.Node{
				Name: "TestField",
				UI:   &node.UIMetadata{MinLen: tt.minLen},
			}
			p := path.NewPath("test")

			err := validateMinLen(n, tt.value, p)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateMinLen() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Field != "TestField" {
				t.Errorf("validateMinLen() error.Field = %s, want TestField", err.Field)
			}
		})
	}
}

func TestValidateMaxLen(t *testing.T) {
	tests := []struct {
		name    string
		maxLen  int
		value   string
		wantErr bool
	}{
		{"no constraint", 0, "very long string", false},
		{"valid", 10, "hello", false},
		{"exact match", 5, "hello", false},
		{"too long", 3, "hello", true},
		{"unicode aware - café", 5, "café", false},
		{"unicode aware - 日本語", 3, "日本語", false},
		{"empty with maxlen 0", 0, "", false},
		{"empty with maxlen 10", 10, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &node.Node{
				Name: "TestField",
				UI:   &node.UIMetadata{MaxLen: tt.maxLen},
			}
			p := path.NewPath("test")

			err := validateMaxLen(n, tt.value, p)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateMaxLen() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Field != "TestField" {
				t.Errorf("validateMaxLen() error.Field = %s, want TestField", err.Field)
			}
		})
	}
}

func TestValidatePattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		value   string
		wantErr bool
	}{
		// No pattern
		{"no pattern", "", "anything", false},

		// Preset patterns
		{"email preset valid", "email", "user@example.com", false},
		{"email preset invalid", "email", "not-an-email", true},
		{"url preset valid http", "url", "http://example.com", false},
		{"url preset valid https", "url", "https://example.com/path", false},
		{"url preset invalid", "url", "not a url", true},
		{"phone preset valid", "phone", "+1 234 567 8900", false},
		{"phone preset valid simple", "phone", "1234567890", false},
		{"phone preset invalid", "phone", "abc", true},

		// Custom regex
		{"custom regex valid", "^[a-z]+$", "hello", false},
		{"custom regex invalid", "^[a-z]+$", "Hello123", true},
		{"custom regex digits", "^[0-9]{3}$", "123", false},
		{"custom regex digits invalid", "^[0-9]{3}$", "12", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &node.Node{
				Name: "TestField",
				UI:   &node.UIMetadata{Pattern: tt.pattern},
			}
			p := path.NewPath("test")

			err := validatePattern(n, tt.value, p)

			if (err != nil) != tt.wantErr {
				t.Errorf("validatePattern() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Field != "TestField" {
				t.Errorf("validatePattern() error.Field = %s, want TestField", err.Field)
			}
		})
	}
}

func TestResolvePattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    string
	}{
		{"email preset", "email", presetPatterns["email"]},
		{"url preset", "url", presetPatterns["url"]},
		{"phone preset", "phone", presetPatterns["phone"]},
		{"custom pattern", "^[a-z]+$", "^[a-z]+$"},
		{"unknown preset", "unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolvePattern(tt.pattern)
			if got != tt.want {
				t.Errorf("resolvePattern(%q) = %q, want %q", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestFormatPatternError(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    string
	}{
		{"email preset", "email", "must be a valid email address"},
		{"url preset", "url", "must be a valid URL"},
		{"phone preset", "phone", "must be a valid phone number"},
		{"custom pattern", "^[a-z]+$", "must match pattern: ^[a-z]+$"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatPatternError(tt.pattern)
			if got != tt.want {
				t.Errorf("formatPatternError(%q) = %q, want %q", tt.pattern, got, tt.want)
			}
		})
	}
}
