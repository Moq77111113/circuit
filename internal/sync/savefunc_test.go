package sync

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSaveFunc_Custom verifies that WithSaveFunc() replaces default persistence.
func TestSaveFunc_Custom(t *testing.T) {
	type Cfg struct {
		Port int `yaml:"port"`
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	customPath := filepath.Join(dir, "custom.yaml")

	err := os.WriteFile(path, []byte("port: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var cfg Cfg
	var customCalled bool

	store, err := Load(Config{
		Path:       path,
		Cfg:        &cfg,
		AutoReload: false,
		Options: []Option{
			WithSaveFunc(func(c any, p string) error {
				customCalled = true
				// Save to custom path instead
				return os.WriteFile(customPath, []byte("port: 7777"), 0644)
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	// Call Save()
	err = store.Save()
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	if !customCalled {
		t.Error("custom SaveFunc should have been called")
	}

	// Verify custom path was written
	if _, err := os.Stat(customPath); os.IsNotExist(err) {
		t.Error("custom path should exist")
	}
}
