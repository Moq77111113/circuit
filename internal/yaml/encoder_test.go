package yaml

import (
	"bytes"
	"testing"
)

func TestEncode_SimpleConfig(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	cfg := Config{Port: 8080}

	data, err := Encode(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("expected non-empty output")
	}

	// Verify it contains the expected value
	if !bytes.Contains(data, []byte("port: 8080")) {
		t.Errorf("expected 'port: 8080' in output, got: %s", data)
	}
}

func TestEncode_MultipleFields(t *testing.T) {
	type Config struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		TLS  bool   `yaml:"tls"`
	}

	cfg := Config{
		Host: "localhost",
		Port: 8080,
		TLS:  true,
	}

	data, err := Encode(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(data, []byte("host: localhost")) {
		t.Error("expected 'host: localhost' in output")
	}
	if !bytes.Contains(data, []byte("port: 8080")) {
		t.Error("expected 'port: 8080' in output")
	}
	if !bytes.Contains(data, []byte("tls: true")) {
		t.Error("expected 'tls: true' in output")
	}
}

func TestEncode_RoundTrip(t *testing.T) {
	type Config struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		TLS  bool   `yaml:"tls"`
	}

	original := Config{
		Host: "example.com",
		Port: 9000,
		TLS:  false,
	}

	// Encode
	data, err := Encode(&original)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	// Parse back
	var decoded Config
	err = Parse(data, &decoded)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	// Verify values match
	if decoded.Host != original.Host {
		t.Errorf("host mismatch: got %s, want %s", decoded.Host, original.Host)
	}
	if decoded.Port != original.Port {
		t.Errorf("port mismatch: got %d, want %d", decoded.Port, original.Port)
	}
	if decoded.TLS != original.TLS {
		t.Errorf("tls mismatch: got %v, want %v", decoded.TLS, original.TLS)
	}
}

func TestEncode_Value(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	cfg := Config{Port: 3000}

	// Should work with value (not just pointer)
	data, err := Encode(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(data, []byte("port: 3000")) {
		t.Errorf("expected 'port: 3000' in output, got: %s", data)
	}
}

func TestEncode_EmptyStruct(t *testing.T) {
	type Empty struct{}

	e := Empty{}
	data, err := Encode(&e)
	if err != nil {
		t.Fatal(err)
	}

	// Empty struct should produce valid (minimal) YAML
	if len(data) == 0 {
		t.Error("expected some output for empty struct")
	}
}
