package sync

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOnError_AutoReloadParseFailure(t *testing.T) {
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

	var receivedError error

	store, err := Load(Config{
		Path:       path,
		Cfg:        &cfg,
		AutoReload: true,
		Options: []Option{
			WithOnError(func(err error) {
				receivedError = err
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	// Write invalid YAML
	err = os.WriteFile(path, []byte("port: invalid"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for watcher to trigger reload
	time.Sleep(200 * time.Millisecond)

	// Should have received parse error
	if receivedError == nil {
		t.Error("expected error callback to be called for parse failure")
	}

	if !errors.Is(receivedError, ErrAutoReloadParse) {
		t.Errorf("expected ErrAutoReloadParse, got %v", receivedError)
	}
}

func TestOnError_ReadFailure(t *testing.T) {
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
		AutoReload: true,
		Options: []Option{
			WithOnError(func(err error) {
				_ = err
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	// Delete file to cause read error
	err = os.Remove(path)
	if err != nil {
		t.Fatal(err)
	}

	// Recreate file to trigger watcher
	err = os.WriteFile(path, []byte("port: 9000"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	// Remove again to cause error
	err = os.Remove(path)
	if err != nil {
		t.Fatal(err)
	}

	// Trigger change event
	err = os.WriteFile(path+"_trigger", []byte("dummy"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)
}

func TestOnError_SaveFuncFailure(t *testing.T) {
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

	saveErr := errors.New("custom save failed")

	store, err := Load(Config{
		Path:       path,
		Cfg:        &cfg,
		AutoReload: false,
		Options: []Option{
			WithSaveFunc(func(cfg any, path string) error {
				return saveErr
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	// Try to save
	err = store.Save()

	// Should receive the custom error
	if err == nil {
		t.Error("expected save to fail with custom error")
	}

	if !errors.Is(err, saveErr) {
		t.Errorf("expected error to wrap saveErr, got %v", err)
	}
}

func TestOnError_NotSet(t *testing.T) {
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

	// No onError callback
	store, err := Load(Config{
		Path:       path,
		Cfg:        &cfg,
		AutoReload: true,
		Options:    []Option{},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	// Write invalid YAML
	err = os.WriteFile(path, []byte("port: invalid"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for watcher
	time.Sleep(200 * time.Millisecond)

	// Should not panic (silent failure is acceptable when callback not set)
}
