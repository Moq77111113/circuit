package json

import (
	"testing"

	"github.com/moq77111113/circuit/internal/codec"
)

func TestCodecParse(t *testing.T) {
	c := Codec{}
	data := []byte(`{"port": 8080, "host": "localhost"}`)

	var cfg testConfig
	err := c.Parse(data, &cfg)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", cfg.Port)
	}
}

func TestCodecEncode(t *testing.T) {
	c := Codec{}
	cfg := testConfig{Port: 8080, Host: "localhost"}

	data, err := c.Encode(cfg)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestCodecRoundTrip(t *testing.T) {
	c := Codec{}
	original := testConfig{Port: 8080, Host: "localhost"}

	data, err := c.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded testConfig
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
