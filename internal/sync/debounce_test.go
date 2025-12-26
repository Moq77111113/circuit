package sync

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestDebounce_FormSubmitIgnored(t *testing.T) {
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

	var mu sync.Mutex
	eventCount := 0
	var lastSource Source

	store, err := Load(Config{
		Path:       path,
		Cfg:        &cfg,
		AutoReload: true,
		Options: []Option{
			WithOnChange(func(e ChangeEvent) {
				mu.Lock()
				eventCount++
				lastSource = e.Source
				mu.Unlock()
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	// Simulate form submit with save
	store.WithLock(func() {
		cfg.Port = 9000
	})

	// Mark form submit (simulates handler behavior)
	store.MarkFormSubmit()

	// Save to disk (this will trigger watcher)
	if err := store.Save(); err != nil {
		t.Fatal(err)
	}

	// Wait for potential watcher event
	time.Sleep(100 * time.Millisecond)

	// Should NOT receive SourceFileChange because of debounce
	mu.Lock()
	count := eventCount
	source := lastSource
	mu.Unlock()

	if count > 0 {
		t.Errorf("expected 0 events (debounced), got %d with source %v", count, source)
	}
}

func TestDebounce_ExternalChangeDetected(t *testing.T) {
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

	var mu sync.Mutex
	eventCount := 0
	var lastSource Source

	store, err := Load(Config{
		Path:       path,
		Cfg:        &cfg,
		AutoReload: true,
		Options: []Option{
			WithOnChange(func(e ChangeEvent) {
				mu.Lock()
				eventCount++
				lastSource = e.Source
				mu.Unlock()
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Stop()

	// Wait for debounce window to pass
	time.Sleep(600 * time.Millisecond)

	// External file change (not from form submit)
	err = os.WriteFile(path, []byte("port: 9000"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for watcher to detect change
	time.Sleep(200 * time.Millisecond)

	// Should receive SourceFileChange
	mu.Lock()
	count := eventCount
	source := lastSource
	mu.Unlock()

	if count != 1 {
		t.Errorf("expected 1 event, got %d", count)
	}

	if source != SourceFileChange {
		t.Errorf("expected SourceFileChange, got %v", source)
	}

	var port int
	store.WithLock(func() {
		port = cfg.Port
	})

	if port != 9000 {
		t.Errorf("expected port 9000 after external change, got %d", port)
	}
}
