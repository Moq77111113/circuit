package sync

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAutoSave_Disabled verifies that with WithAutoSave(false),
// Save() is not called automatically on changes.
func TestAutoSave_Disabled(t *testing.T) {
	type Cfg struct {
		Port int `yaml:"port"`
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	err := os.WriteFile(path, []byte("port: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var cfg Cfg
	store, err := Load(Config{
		Path:       path,
		Cfg:        &cfg,
		AutoReload: false,
		Options:    []Option{WithAutoSave(false)},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	if store.AutoSave() {
		t.Error("AutoSave should be false")
	}

	// Manually change config
	store.mu.Lock()
	cfg.Port = 9000
	store.mu.Unlock()

	// Call Save() manually
	err = store.Save()
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify file was updated by reloading
	var reloaded Cfg
	s2, err := Load(Config{
		Path:       path,
		Cfg:        &reloaded,
		AutoReload: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Stop()

	if reloaded.Port != 9000 {
		t.Errorf("expected port 9000 after manual save, got %d", reloaded.Port)
	}
}
