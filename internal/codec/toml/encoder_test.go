package toml

import (
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	cfg := testConfig{Port: 8080, Host: "localhost"}

	data, err := Encode(cfg)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	result := string(data)
	if !strings.Contains(result, "port = 8080") {
		t.Errorf("expected port = 8080 in output, got: %s", result)
	}
	if !strings.Contains(result, `host = "localhost"`) {
		t.Errorf("expected host = \"localhost\" in output, got: %s", result)
	}
}

func TestEncodeMultipleFields(t *testing.T) {
	type Config struct {
		Port    int    `toml:"port"`
		Host    string `toml:"host"`
		Enabled bool   `toml:"enabled"`
	}

	cfg := Config{Port: 8080, Host: "localhost", Enabled: true}

	data, err := Encode(cfg)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	result := string(data)
	if !strings.Contains(result, "port = 8080") {
		t.Error("missing port field")
	}
	if !strings.Contains(result, `host = "localhost"`) {
		t.Error("missing host field")
	}
	if !strings.Contains(result, "enabled = true") {
		t.Error("missing enabled field")
	}
}

func TestEncodeRoundTrip(t *testing.T) {
	original := testConfig{Port: 8080, Host: "localhost"}

	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded testConfig
	err = Parse(data, &decoded)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if decoded != original {
		t.Errorf("round-trip failed: expected %+v, got %+v", original, decoded)
	}
}

func TestEncodeValue(t *testing.T) {
	cfg := testConfig{Port: 3000, Host: "example.com"}

	data, err := Encode(cfg)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestEncodeEmpty(t *testing.T) {
	cfg := testConfig{}

	data, err := Encode(cfg)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty output for empty struct")
	}
}
