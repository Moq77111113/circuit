package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/sync"
)

// TestHandler_Confirm verifies that action=confirm applies changes.
func TestHandler_Confirm(t *testing.T) {
	type Cfg struct {
		Port int `circuit:"Port,number" yaml:"port"`
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	err := os.WriteFile(path, []byte("port: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var cfg Cfg
	s, err := ast.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	store, err := sync.Load(sync.Config{
		Path:       path,
		Cfg:        &cfg,
		AutoReload: false,
		Options: []sync.Option{
			sync.WithAutoApply(false), // Preview mode
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	h := New(Config{
		Schema: s,
		Cfg:    &cfg,
		Path:   path,
		Title:  "Test Confirm",
		Brand:  true,
		Store:  store,
	})

	// Submit with action=confirm
	form := url.Values{}
	form.Set("Port", "9000")
	form.Set("action", "confirm")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	// Verify config WAS updated (confirm applied changes)
	var port int
	store.WithLock(func() {
		port = cfg.Port
	})

	if port != 9000 {
		t.Errorf("config should be updated after confirm, expected 9000 got %d", port)
	}

	// Verify redirect happened
	if w.Code != http.StatusSeeOther {
		t.Errorf("expected redirect 303 after confirm, got %d", w.Code)
	}
}
