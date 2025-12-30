package yaml

import "testing"

func TestParse_SimpleConfig(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	input := []byte("port: 8080")
	var cfg Config

	err := Parse(input, &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
}

func TestParse_MultipleFields(t *testing.T) {
	type Config struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		TLS  bool   `yaml:"tls"`
	}

	input := []byte(`
host: localhost
port: 8080
tls: true
`)

	var cfg Config
	err := Parse(input, &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if !cfg.TLS {
		t.Error("expected tls true, got false")
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	input := []byte("port: [invalid")
	var cfg Config

	err := Parse(input, &cfg)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestParse_NonPointer(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	input := []byte("port: 8080")
	var cfg Config

	err := Parse(input, cfg)
	if err == nil {
		t.Fatal("expected error for non-pointer destination")
	}
}

func TestParse_EmptyInput(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	input := []byte("")
	var cfg Config

	err := Parse(input, &cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Empty YAML should leave fields at zero values
	if cfg.Port != 0 {
		t.Errorf("expected port 0, got %d", cfg.Port)
	}
}
