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

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/auth"
	"github.com/moq77111113/circuit/internal/sync"
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

	store, err := sync.Load(sync.Config{Path: path, Cfg: &cfg, AutoReload: true})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	h := New(Config{
		Schema: s,
		Cfg:    &cfg,
		Path:   path,
		Title:  "Test",
		Brand:  true,
		Store:  store,
	})

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

	store, err := sync.Load(sync.Config{Path: path, Cfg: &cfg, AutoReload: true})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	h := New(Config{
		Schema: s,
		Cfg:    &cfg,
		Path:   path,
		Title:  "Test",
		Brand:  true,
		Store:  store,
	})

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

	store, err := sync.Load(sync.Config{Path: path, Cfg: &cfg, AutoReload: true})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	h := New(Config{
		Schema: s,
		Cfg:    &cfg,
		Path:   path,
		Title:  "Test",
		Brand:  true,
		Store:  store,
	})

	req := httptest.NewRequest("DELETE", "/", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rec.Code)
	}
}

func TestHandler_GET_WithValidAuth(t *testing.T) {
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

	store, err := sync.Load(sync.Config{Path: path, Cfg: &cfg, AutoReload: true})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	authenticator := auth.Basic{
		Username: "admin",
		Password: "secret",
	}

	h := New(Config{
		Schema:        s,
		Cfg:           &cfg,
		Path:          path,
		Title:         "Test",
		Brand:         true,
		Store:         store,
		Authenticator: authenticator,
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("admin", "secret")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestHandler_GET_WithInvalidAuth(t *testing.T) {
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

	store, err := sync.Load(sync.Config{Path: path, Cfg: &cfg, AutoReload: true})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	authenticator := auth.Basic{
		Username: "admin",
		Password: "secret",
	}

	h := New(Config{
		Schema:        s,
		Cfg:           &cfg,
		Path:          path,
		Title:         "Test",
		Brand:         true,
		Store:         store,
		Authenticator: authenticator,
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("admin", "wrongpassword")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}

	wwwAuth := rec.Header().Get("WWW-Authenticate")
	if wwwAuth != `Basic realm="Circuit"` {
		t.Errorf("expected WWW-Authenticate header, got: %q", wwwAuth)
	}
}

func TestHandler_POST_WithValidAuth(t *testing.T) {
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

	store, err := sync.Load(sync.Config{Path: path, Cfg: &cfg, AutoReload: true})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	authenticator := auth.Basic{
		Username: "user",
		Password: "pass",
	}

	h := New(Config{
		Schema:        s,
		Cfg:           &cfg,
		Path:          path,
		Title:         "Test",
		Brand:         true,
		Store:         store,
		Authenticator: authenticator,
	})

	form := url.Values{}
	form.Set("Host", "example.com")

	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("user", "pass")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}
}

func TestHandler_POST_WithInvalidAuth(t *testing.T) {
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

	store, err := sync.Load(sync.Config{Path: path, Cfg: &cfg, AutoReload: true})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	authenticator := auth.Basic{
		Username: "user",
		Password: "pass",
	}

	h := New(Config{
		Schema:        s,
		Cfg:           &cfg,
		Path:          path,
		Title:         "Test",
		Brand:         true,
		Store:         store,
		Authenticator: authenticator,
	})

	form := url.Values{}
	form.Set("Host", "hacked.com")

	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("attacker", "wrongpass")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestHandler_NoAuthConfigured(t *testing.T) {
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

	store, err := sync.Load(sync.Config{Path: path, Cfg: &cfg, AutoReload: true})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	h := New(Config{
		Schema:        s,
		Cfg:           &cfg,
		Path:          path,
		Title:         "Test",
		Brand:         true,
		Store:         store,
		Authenticator: nil, // defaults to None
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 when no auth configured, got %d", rec.Code)
	}
}

func TestHandler_ForwardAuth(t *testing.T) {
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

	store, err := sync.Load(sync.Config{Path: path, Cfg: &cfg, AutoReload: true})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	authenticator := auth.Forward{
		SubjectHeader: "X-Forwarded-User",
	}

	h := New(Config{
		Schema:        s,
		Cfg:           &cfg,
		Path:          path,
		Title:         "Test",
		Brand:         true,
		Store:         store,
		Authenticator: authenticator,
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-User", "alice@example.com")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
