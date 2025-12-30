package json

import "testing"

type testConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

func TestParse(t *testing.T) {
	data := []byte(`{"port": 8080, "host": "localhost"}`)

	var cfg testConfig
	err := Parse(data, &cfg)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", cfg.Port)
	}
	if cfg.Host != "localhost" {
		t.Errorf("expected Host=localhost, got %q", cfg.Host)
	}
}

func TestParseNested(t *testing.T) {
	type Database struct {
		Port int    `json:"port"`
		Host string `json:"host"`
	}
	type Config struct {
		DB Database `json:"database"`
	}

	data := []byte(`{"database": {"port": 5432, "host": "postgres"}}`)

	var cfg Config
	err := Parse(data, &cfg)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.DB.Port != 5432 {
		t.Errorf("expected DB.Port=5432, got %d", cfg.DB.Port)
	}
	if cfg.DB.Host != "postgres" {
		t.Errorf("expected DB.Host=postgres, got %q", cfg.DB.Host)
	}
}

func TestParseArray(t *testing.T) {
	type Config struct {
		Servers []string `json:"servers"`
	}

	data := []byte(`{"servers": ["server1", "server2", "server3"]}`)

	var cfg Config
	err := Parse(data, &cfg)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(cfg.Servers) != 3 {
		t.Errorf("expected 3 servers, got %d", len(cfg.Servers))
	}
}

func TestParseInvalidJSON(t *testing.T) {
	data := []byte(`{invalid json`)

	var cfg testConfig
	err := Parse(data, &cfg)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParseNonPointer(t *testing.T) {
	data := []byte(`{"port": 8080}`)

	var cfg testConfig
	err := Parse(data, cfg)
	if err == nil {
		t.Fatal("expected error for non-pointer, got nil")
	}
}

func TestParseEmpty(t *testing.T) {
	data := []byte(`{}`)

	var cfg testConfig
	err := Parse(data, &cfg)
	if err != nil {
		t.Errorf("unexpected error for empty JSON: %v", err)
	}
}
