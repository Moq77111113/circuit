package form

import (
	"net/url"
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
)

func TestSliceEdit_OneItemOnly(t *testing.T) {
	type Service struct {
		Name string `yaml:"name" circuit:"type:text"`
		Port int    `yaml:"port" circuit:"type:number"`
	}

	type Config struct {
		Services []Service `yaml:"services"`
	}

	cfg := &Config{
		Services: []Service{
			{Name: "service1", Port: 8080},
			{Name: "service2", Port: 8081},
			{Name: "service3", Port: 8082},
		},
	}

	schema, err := ast.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Set("Services.0.Name", "updated1")
	form.Set("Services.0.Port", "9000")

	err = Apply(cfg, schema, form)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Services) != 3 {
		t.Fatalf("expected 3 services, got %d", len(cfg.Services))
	}

	if cfg.Services[0].Name != "updated1" {
		t.Errorf("expected Services[0].Name=updated1, got %s", cfg.Services[0].Name)
	}
	if cfg.Services[0].Port != 9000 {
		t.Errorf("expected Services[0].Port=9000, got %d", cfg.Services[0].Port)
	}

	if cfg.Services[1].Name != "service2" {
		t.Errorf("expected Services[1].Name=service2, got %s (PRESERVED)", cfg.Services[1].Name)
	}
	if cfg.Services[1].Port != 8081 {
		t.Errorf("expected Services[1].Port=8081, got %d (PRESERVED)", cfg.Services[1].Port)
	}

	if cfg.Services[2].Name != "service3" {
		t.Errorf("expected Services[2].Name=service3, got %s (PRESERVED)", cfg.Services[2].Name)
	}
	if cfg.Services[2].Port != 8082 {
		t.Errorf("expected Services[2].Port=8082, got %d (PRESERVED)", cfg.Services[2].Port)
	}
}

func TestSliceEdit_FromRootView(t *testing.T) {
	type Service struct {
		Name string `yaml:"name" circuit:"type:text"`
		Port int    `yaml:"port" circuit:"type:number"`
	}

	type Config struct {
		Host     string    `yaml:"host" circuit:"type:text"`
		Services []Service `yaml:"services"`
	}

	cfg := &Config{
		Host: "localhost",
		Services: []Service{
			{Name: "service1", Port: 8080},
			{Name: "service2", Port: 8081},
		},
	}

	schema, err := ast.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Set("Host", "example.com")
	form.Set("Services.0.Name", "api")
	form.Set("Services.0.Port", "3000")
	form.Set("Services.1.Name", "service2")
	form.Set("Services.1.Port", "8081")

	err = Apply(cfg, schema, form)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Host != "example.com" {
		t.Errorf("expected Host=example.com, got %s", cfg.Host)
	}

	if len(cfg.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(cfg.Services))
	}

	if cfg.Services[0].Name != "api" {
		t.Errorf("expected Services[0].Name=api, got '%s'", cfg.Services[0].Name)
	}
	if cfg.Services[0].Port != 3000 {
		t.Errorf("expected Services[0].Port=3000, got %d", cfg.Services[0].Port)
	}

	if cfg.Services[1].Name != "service2" {
		t.Errorf("expected Services[1].Name=service2, got '%s'", cfg.Services[1].Name)
	}
	if cfg.Services[1].Port != 8081 {
		t.Errorf("expected Services[1].Port=8081, got %d", cfg.Services[1].Port)
	}
}
