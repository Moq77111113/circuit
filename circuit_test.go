package circuit

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type TestConfig struct {
	Host string `yaml:"host" circuit:"text,help:Server hostname"`
	Port int    `yaml:"port" circuit:"number,help:Server port,required"`
	TLS  bool   `yaml:"tls" circuit:"checkbox,help:Enable TLS"`
}

func TestUI_GET(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{Host: "localhost", Port: 8080, TLS: true}
	err := os.WriteFile(path, []byte("host: localhost\nport: 8080\ntls: true"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	handler, err := From(&cfg, WithPath(path))
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "localhost") {
		t.Error("expected form to contain current host value")
	}
	if !strings.Contains(body, "8080") {
		t.Error("expected form to contain current port value")
	}
	if !strings.Contains(body, "Server hostname") {
		t.Error("expected form to contain help text")
	}
}

func TestUI_POST(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{Host: "localhost", Port: 8080}
	err := os.WriteFile(path, []byte("host: localhost\nport: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var callbackCalled atomic.Bool
	h, err := From(&cfg, WithPath(path), WithOnChange(func(e ChangeEvent) {
		callbackCalled.Store(true)
	}))
	if err != nil {
		t.Fatal(err)
	}

	// Submit new values
	form := url.Values{}
	form.Set("Host", "example.com")
	form.Set("Port", "9000")
	form.Set("TLS", "true")

	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}

	// Wait for reload to propagate
	time.Sleep(200 * time.Millisecond)

	// Read file to verify it was saved
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
	if !strings.Contains(string(saved), "true") {
		t.Errorf("expected saved config to contain true for TLS, got: %s", saved)
	}

	if !callbackCalled.Load() {
		t.Error("expected callback to be called")
	}
}

func TestUI_NoPath(t *testing.T) {
	cfg := TestConfig{}
	_, err := From(&cfg)
	if err == nil {
		t.Fatal("expected error when path is not provided")
	}
}

func TestUI_InvalidPath(t *testing.T) {
	cfg := TestConfig{}
	_, err := From(&cfg, WithPath("/nonexistent/config.yaml"))
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestUI_WithTitle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
	err := os.WriteFile(path, []byte("host: localhost"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	handler, err := From(&cfg, WithPath(path), WithTitle("My Settings"))
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "My Settings") {
		t.Error("expected custom title in response")
	}
}

func TestUI_NonPointer(t *testing.T) {
	cfg := TestConfig{}
	_, err := From(cfg, WithPath("/tmp/config.yaml"))
	if err == nil {
		t.Fatal("expected error for non-pointer config")
	}
}

func TestUI_WithBasicAuth(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
	err := os.WriteFile(path, []byte("host: localhost\nport: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	auth := NewBasicAuth("admin", "secret123")
	handler, err := From(&cfg, WithPath(path), WithAuth(auth))
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("admin", "secret123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 with valid auth, got %d", rec.Code)
	}
}

func TestUI_WithBasicAuthUnauthorized(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
	err := os.WriteFile(path, []byte("host: localhost"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	auth := NewBasicAuth("admin", "secret")
	handler, err := From(&cfg, WithPath(path), WithAuth(auth))
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	// No auth header
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 without auth, got %d", rec.Code)
	}
}

func TestUI_WithForwardAuth(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
	err := os.WriteFile(path, []byte("host: localhost"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	auth := NewForwardAuth("X-Forwarded-User", map[string]string{
		"email": "X-Forwarded-Email",
	})
	handler, err := From(&cfg, WithPath(path), WithAuth(auth))
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-User", "alice@example.com")
	req.Header.Set("X-Forwarded-Email", "alice@example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 with forward auth, got %d", rec.Code)
	}
}

func TestUI_NoAuthConfigured(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
	err := os.WriteFile(path, []byte("host: localhost"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	handler, err := From(&cfg, WithPath(path))
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	// No auth headers
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 when no auth configured, got %d", rec.Code)
	}
}
