package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/moq77111113/circuit/internal/reload"
	"github.com/moq77111113/circuit/internal/ast"
)

type TestConfig struct {
	Host string `yaml:"host" circuit:"type:text"`
	Port int    `yaml:"port" circuit:"type:number"`
	TLS  bool   `yaml:"tls" circuit:"type:checkbox"`
}

func TestHandler_GET(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{Host: "localhost", Port: 8080}
	err := os.WriteFile(path, []byte("host: localhost\nport: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	s, err := ast.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	loader, err := reload.Load(path, &cfg, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer loader.Stop()

	h := New(s, &cfg, path, "Test", true, loader)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "localhost") {
		t.Error("expected response to contain localhost")
	}
	if !strings.Contains(body, "8080") {
		t.Error("expected response to contain 8080")
	}
}

func TestHandler_POST(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{Host: "localhost", Port: 8080}
	err := os.WriteFile(path, []byte("host: localhost\nport: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	s, err := ast.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	loader, err := reload.Load(path, &cfg, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer loader.Stop()

	h := New(s, &cfg, path, "Test", true, loader)

	form := url.Values{}
	form.Set("Host", "example.com")
	form.Set("Port", "9000")
	form.Set("TLS", "on")

	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}

	// Wait for reload
	time.Sleep(100 * time.Millisecond)

	// Verify file was updated
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(saved), "example.com") {
		t.Errorf("expected saved config to contain example.com, got: %s", saved)
	}
	if !strings.Contains(string(saved), "9000") {
		t.Errorf("expected saved config to contain 9000, got: %s", saved)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
	err := os.WriteFile(path, []byte("host: localhost"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	s, err := ast.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	loader, err := reload.Load(path, &cfg, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer loader.Stop()

	h := New(s, &cfg, path, "Test", true, loader)

	req := httptest.NewRequest("DELETE", "/", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rec.Code)
	}
}
