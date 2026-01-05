package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/moq77111113/circuit/internal/actions"
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

func TestHandler_ExecuteAction_Integration(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
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

	executed := false
	h := New(Config{
		Schema: s,
		Cfg:    &cfg,
		Path:   path,
		Title:  "Test",
		Brand:  true,
		Store:  store,
		Actions: []actions.Def{
			{
				Name:  "test-action",
				Label: "Test Action",
				Run: func(ctx context.Context) error {
					executed = true
					return nil
				},
			},
		},
	})

	form := url.Values{}
	form.Set("action", "execute:test-action")
	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if !executed {
		t.Fatal("expected action to be executed")
	}

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect status %d, got %d", http.StatusSeeOther, rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "/" {
		t.Errorf("expected redirect to /, got %s", location)
	}
}

func TestHandler_ExecuteAction_ErrorDisplay(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
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

	expectedErr := errors.New("action execution failed")
	h := New(Config{
		Schema: s,
		Cfg:    &cfg,
		Path:   path,
		Title:  "Test",
		Brand:  true,
		Store:  store,
		Actions: []actions.Def{
			{
				Name:  "failing-action",
				Label: "Failing Action",
				Run: func(ctx context.Context) error {
					return expectedErr
				},
			},
		},
	})

	form := url.Values{}
	form.Set("action", "execute:failing-action")
	postReq := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postRec := httptest.NewRecorder()

	h.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect status %d, got %d", http.StatusSeeOther, postRec.Code)
	}

	location := postRec.Header().Get("Location")
	if !strings.Contains(location, "?error=") {
		t.Fatalf("expected redirect with error param, got %s", location)
	}

	getReq := httptest.NewRequest("GET", location, nil)
	getRec := httptest.NewRecorder()

	h.ServeHTTP(getRec, getReq)

	body := getRec.Body.String()
	if !strings.Contains(body, "Error:") {
		t.Error("expected error banner in response body")
	}
	if !strings.Contains(body, expectedErr.Error()) {
		t.Errorf("expected error message %q in response body", expectedErr.Error())
	}
}

func TestHandler_HTTPBasePath(t *testing.T) {
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

	// Test that GET request renders correct base path in breadcrumb and sidebar
	req := httptest.NewRequest("GET", "/admin/config", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	// Check that breadcrumb root link uses the correct base path
	if !strings.Contains(body, `href="/admin/config"`) {
		t.Error("expected breadcrumb root link to use /admin/config base path")
	}

	// Check that sidebar Config link uses the correct base path
	if !strings.Contains(body, `<a class="nav__link" href="/admin/config">`) {
		t.Error("expected sidebar Config link to use /admin/config base path")
	}
}

func TestHandler_HTTPBasePath_ActionRedirect(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
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
		Actions: []actions.Def{
			{
				Name:  "test-action",
				Label: "Test Action",
				Run: func(ctx context.Context) error {
					return nil
				},
			},
		},
	})

	// Test that action execution redirects to correct base path
	form := url.Values{}
	form.Set("action", "execute:test-action")
	req := httptest.NewRequest("POST", "/admin/config", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect status %d, got %d", http.StatusSeeOther, rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "/admin/config" {
		t.Errorf("expected redirect to /admin/config, got %s", location)
	}
}

func TestHandler_HTTPBasePath_ActionError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
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

	expectedErr := errors.New("something went wrong")
	h := New(Config{
		Schema: s,
		Cfg:    &cfg,
		Path:   path,
		Title:  "Test",
		Brand:  true,
		Store:  store,
		Actions: []actions.Def{
			{
				Name:  "failing-action",
				Label: "Failing Action",
				Run: func(ctx context.Context) error {
					return expectedErr
				},
			},
		},
	})

	// Test that action error redirects to correct base path with error param
	form := url.Values{}
	form.Set("action", "execute:failing-action")
	req := httptest.NewRequest("POST", "/admin/config", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect status %d, got %d", http.StatusSeeOther, rec.Code)
	}

	location := rec.Header().Get("Location")
	if !strings.HasPrefix(location, "/admin/config?error=") {
		t.Errorf("expected redirect to /admin/config?error=..., got %s", location)
	}
}

func TestHandler_ActionsNotShownInReadOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := TestConfig{}
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
		Schema:   s,
		Cfg:      &cfg,
		Path:     path,
		Title:    "Test",
		Brand:    true,
		ReadOnly: true,
		Store:    store,
		Actions: []actions.Def{
			{
				Name:  "test-action",
				Label: "Test Action",
				Run: func(ctx context.Context) error {
					return nil
				},
			},
		},
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	body := rec.Body.String()
	if strings.Contains(body, `<div class="actions-section">`) {
		t.Error("expected actions section HTML to be absent in read-only mode")
	}
	if strings.Contains(body, "Test Action") {
		t.Error("expected action button to be absent in read-only mode")
	}
}
