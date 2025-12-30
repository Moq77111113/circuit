package yaml

import (
	"testing"

	"github.com/moq77111113/circuit/internal/codec"
)

func TestCodecParse(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	c := Codec{}
	data := []byte("port: 8080")
	var cfg Config

	err := c.Parse(data, &cfg)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", cfg.Port)
	}
}

func TestCodecEncode(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	c := Codec{}
	cfg := Config{Port: 8080}

	data, err := c.Encode(cfg)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	expected := "port: 8080\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestCodecRoundTrip(t *testing.T) {
	type Config struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	}

	c := Codec{}
	original := Config{Port: 8080, Host: "localhost"}

	// Encode
	data, err := c.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode
	var decoded Config
	err = c.Parse(data, &decoded)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if decoded != original {
		t.Errorf("round-trip failed: expected %+v, got %+v", original, decoded)
	}
}

func TestCodecImplementsInterface(t *testing.T) {
	var _ codec.Codec = Codec{}
}
