package codec_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/codec"
	_ "github.com/moq77111113/circuit/internal/codec/json"
	_ "github.com/moq77111113/circuit/internal/codec/toml"
	_ "github.com/moq77111113/circuit/internal/codec/yaml"
)

type e2eConfig struct {
	Port int    `yaml:"port" toml:"port" json:"port"`
	Host string `yaml:"host" toml:"host" json:"host"`
}

func TestE2E_YAMLFormat(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	initialData := []byte("port: 8080\nhost: localhost\n")
	if err := os.WriteFile(cfgPath, initialData, 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Detect and parse
	cdc, err := codec.Detect(cfgPath)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	data, _ := os.ReadFile(cfgPath)
	cfg := e2eConfig{}
	if err := cdc.Parse(data, &cfg); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.Port != 8080 || cfg.Host != "localhost" {
		t.Errorf("parse failed: got port=%d host=%s", cfg.Port, cfg.Host)
	}

	// Modify and encode
	cfg.Port = 9000
	encoded, err := cdc.Encode(cfg)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Write back
	if err := os.WriteFile(cfgPath, encoded, 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Verify YAML format preserved
	saved, _ := os.ReadFile(cfgPath)
	if !strings.Contains(string(saved), "port: 9000") {
		t.Error("YAML format not preserved")
	}
}

func TestE2E_TOMLFormat(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.toml")

	initialData := []byte("port = 8080\nhost = \"localhost\"\n")
	if err := os.WriteFile(cfgPath, initialData, 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	cdc, err := codec.Detect(cfgPath)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	data, _ := os.ReadFile(cfgPath)
	cfg := e2eConfig{}
	if err := cdc.Parse(data, &cfg); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.Port != 8080 || cfg.Host != "localhost" {
		t.Errorf("parse failed: got port=%d host=%s", cfg.Port, cfg.Host)
	}

	cfg.Port = 9000
	encoded, err := cdc.Encode(cfg)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if err := os.WriteFile(cfgPath, encoded, 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	saved, _ := os.ReadFile(cfgPath)
	if !strings.Contains(string(saved), "port = 9000") {
		t.Error("TOML format not preserved")
	}
}

func TestE2E_JSONFormat(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")

	initialData := []byte(`{"port": 8080, "host": "localhost"}`)
	if err := os.WriteFile(cfgPath, initialData, 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	cdc, err := codec.Detect(cfgPath)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	data, _ := os.ReadFile(cfgPath)
	cfg := e2eConfig{}
	if err := cdc.Parse(data, &cfg); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.Port != 8080 || cfg.Host != "localhost" {
		t.Errorf("parse failed: got port=%d host=%s", cfg.Port, cfg.Host)
	}

	cfg.Port = 9000
	encoded, err := cdc.Encode(cfg)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if err := os.WriteFile(cfgPath, encoded, 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	saved, _ := os.ReadFile(cfgPath)
	if !strings.Contains(string(saved), `"port": 9000`) {
		t.Error("JSON format not preserved")
	}
}

func TestE2E_FormatPreservation(t *testing.T) {
	tests := []struct {
		name        string
		ext         string
		initialData string
		checkString string
	}{
		{
			name:        "yaml",
			ext:         ".yaml",
			initialData: "port: 8080\nhost: localhost\n",
			checkString: "port:",
		},
		{
			name:        "yml",
			ext:         ".yml",
			initialData: "port: 8080\nhost: localhost\n",
			checkString: "port:",
		},
		{
			name:        "toml",
			ext:         ".toml",
			initialData: "port = 8080\nhost = \"localhost\"\n",
			checkString: "port =",
		},
		{
			name:        "json",
			ext:         ".json",
			initialData: `{"port": 8080, "host": "localhost"}`,
			checkString: `"port"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			cfgPath := filepath.Join(tmpDir, "config"+tt.ext)

			if err := os.WriteFile(cfgPath, []byte(tt.initialData), 0644); err != nil {
				t.Fatalf("write failed: %v", err)
			}

			cdc, err := codec.Detect(cfgPath)
			if err != nil {
				t.Fatalf("Detect failed: %v", err)
			}

			data, _ := os.ReadFile(cfgPath)
			cfg := e2eConfig{}
			if err := cdc.Parse(data, &cfg); err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			cfg.Port = 9999
			encoded, err := cdc.Encode(cfg)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			if err := os.WriteFile(cfgPath, encoded, 0644); err != nil {
				t.Fatalf("write failed: %v", err)
			}

			saved, _ := os.ReadFile(cfgPath)
			if !strings.Contains(string(saved), tt.checkString) {
				t.Errorf("format not preserved, expected %q in:\n%s", tt.checkString, string(saved))
			}
		})
	}
}

func TestE2E_UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.xml")

	if err := os.WriteFile(cfgPath, []byte("<config></config>"), 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	_, err := codec.Detect(cfgPath)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}

	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("expected 'unsupported format' error, got: %v", err)
	}
}
