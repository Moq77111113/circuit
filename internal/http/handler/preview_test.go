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

// TestHandler_PreviewMode_AutoApplyFalse verifies preview mode behavior.
func TestHandler_PreviewMode_AutoApplyFalse(t *testing.T) {
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
			sync.WithAutoApply(false),
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
		Title:  "Test Preview",
		Brand:  true,
		Store:  store,
	})

	form := url.Values{}
	form.Set("Port", "9000")
	form.Set("action", "save")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	// Verify config in memory was NOT updated
	var port int
	store.WithLock(func() {
		port = cfg.Port
	})

	if port != 8080 {
		t.Errorf("config should NOT be updated in preview mode, expected 8080 got %d", port)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 OK for preview, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "9000") {
		t.Error("preview should show submitted value 9000")
	}
	if !strings.Contains(body, "Confirm") || !strings.Contains(body, "Cancel") {
		t.Error("preview should have Confirm and Cancel buttons")
	}
}
