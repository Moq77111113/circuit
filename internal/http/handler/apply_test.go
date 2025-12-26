package handler

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/sync"
)

// TestHandler_Apply verifies that Handler.Apply() manually applies form data.
// This is used in preview mode to confirm changes after user review.
func TestHandler_Apply(t *testing.T) {
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
			sync.WithAutoSave(true),   // Auto-save when applied
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
		Title:  "Test Apply",
		Brand:  true,
		Store:  store,
	})

	// Prepare form data
	formData := url.Values{}
	formData.Set("Port", "9000")

	// Call Apply() to manually apply changes
	err = h.Apply(formData)
	if err != nil {
		t.Fatalf("Apply() failed: %v", err)
	}

	// Verify config was updated in memory
	var port int
	store.WithLock(func() {
		port = cfg.Port
	})

	if port != 9000 {
		t.Errorf("expected port 9000 after Apply(), got %d", port)
	}

	// Verify file was saved (autoSave=true)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "9000") {
		t.Errorf("file should contain 9000, got: %s", string(data))
	}
}
