package form

import (
	"net/url"
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
)

type ServerConfig struct {
	Host string `yaml:"host" circuit:"type:text"`
	Port int    `yaml:"port" circuit:"type:number"`
}

type DatabaseConfig struct {
	Host string `yaml:"host" circuit:"type:text"`
	Port int    `yaml:"port" circuit:"type:number"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func TestApplyForm_FocusedSection_PreservesOtherFields(t *testing.T) {
	cfg := &AppConfig{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host: "db.example.com",
			Port: 5432,
		},
	}

	schema, err := ast.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate form submission with ONLY Server fields (focused section)
	form := url.Values{}
	form.Set("Server.Host", "newhost")
	form.Set("Server.Port", "9000")

	err = Apply(cfg, schema, form)
	if err != nil {
		t.Fatal(err)
	}

	// Server should be updated
	if cfg.Server.Host != "newhost" {
		t.Errorf("expected Server.Host=newhost, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 9000 {
		t.Errorf("expected Server.Port=9000, got %d", cfg.Server.Port)
	}

	// Database should be PRESERVED (not cleared)
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("expected Database.Host preserved, got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected Database.Port=5432, got %d", cfg.Database.Port)
	}
}

func TestApplyForm_EmptyStringWhenPresent_ClearsField(t *testing.T) {
	type TestConfig struct {
		Name string `yaml:"name" circuit:"type:text"`
	}

	cfg := &TestConfig{Name: "oldvalue"}
	schema, err := ast.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Field IS in form but with empty value
	form := url.Values{}
	form.Set("Name", "")

	err = Apply(cfg, schema, form)
	if err != nil {
		t.Fatal(err)
	}

	// Field should be cleared
	if cfg.Name != "" {
		t.Errorf("expected Name cleared, got %s", cfg.Name)
	}
}

func TestApplyForm_CheckboxAbsent_PreservesValue(t *testing.T) {
	type TestConfig struct {
		Debug bool `yaml:"debug" circuit:"type:checkbox"`
		Port  int  `yaml:"port" circuit:"type:number"`
	}

	cfg := &TestConfig{Debug: true, Port: 8080}
	schema, err := ast.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Form contains Port but NOT Debug (filtered out)
	form := url.Values{}
	form.Set("Port", "9000")

	err = Apply(cfg, schema, form)
	if err != nil {
		t.Fatal(err)
	}

	// Debug should be preserved
	if !cfg.Debug {
		t.Error("expected Debug=true preserved, got false")
	}
	if cfg.Port != 9000 {
		t.Errorf("expected Port=9000, got %d", cfg.Port)
	}
}

func TestApplyForm_SliceAbsent_PreservesSlice(t *testing.T) {
	type TestConfig struct {
		Items []string `yaml:"items" circuit:"type:text"`
		Port  int      `yaml:"port" circuit:"type:number"`
	}

	cfg := &TestConfig{
		Items: []string{"a", "b", "c"},
		Port:  8080,
	}
	schema, err := ast.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Form contains Port but NOT Items (filtered out)
	form := url.Values{}
	form.Set("Port", "9000")

	err = Apply(cfg, schema, form)
	if err != nil {
		t.Fatal(err)
	}

	// Items should be preserved
	if len(cfg.Items) != 3 {
		t.Errorf("expected Items preserved with 3 elements, got %d", len(cfg.Items))
	}
	if cfg.Port != 9000 {
		t.Errorf("expected Port=9000, got %d", cfg.Port)
	}
}

func TestApplyForm_NestedStructPartial_PreservesUnsubmitted(t *testing.T) {
	type TLSConfig struct {
		Cert string `yaml:"cert" circuit:"type:text"`
		Key  string `yaml:"key" circuit:"type:text"`
	}

	type ServerConfig struct {
		Host string    `yaml:"host" circuit:"type:text"`
		Port int       `yaml:"port" circuit:"type:number"`
		TLS  TLSConfig `yaml:"tls"`
	}

	type TestConfig struct {
		Server ServerConfig `yaml:"server"`
	}

	cfg := &TestConfig{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
			TLS: TLSConfig{
				Cert: "cert.pem",
				Key:  "key.pem",
			},
		},
	}

	schema, err := ast.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Submit ONLY Server.Host
	form := url.Values{}
	form.Set("Server.Host", "newhost")

	err = Apply(cfg, schema, form)
	if err != nil {
		t.Fatal(err)
	}

	// Server.Host should be updated
	if cfg.Server.Host != "newhost" {
		t.Errorf("expected Server.Host=newhost, got %s", cfg.Server.Host)
	}

	// Server.Port should be preserved
	if cfg.Server.Port != 8080 {
		t.Errorf("expected Server.Port=8080, got %d", cfg.Server.Port)
	}

	// Server.TLS should be preserved
	if cfg.Server.TLS.Cert != "cert.pem" {
		t.Errorf("expected TLS.Cert preserved, got %s", cfg.Server.TLS.Cert)
	}
	if cfg.Server.TLS.Key != "key.pem" {
		t.Errorf("expected TLS.Key preserved, got %s", cfg.Server.TLS.Key)
	}
}
