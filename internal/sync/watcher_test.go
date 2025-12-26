package sync

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatch_FileModified(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	// Create initial file
	err := os.WriteFile(path, []byte("port: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var called atomic.Bool
	callback := func() {
		called.Store(true)
	}

	w, err := Watch(path, callback, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify file
	err = os.WriteFile(path, []byte("port: 9000"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for callback
	time.Sleep(200 * time.Millisecond)

	if !called.Load() {
		t.Error("callback was not called after file modification")
	}
}

func TestWatch_Stop(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	err := os.WriteFile(path, []byte("port: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var called atomic.Bool
	callback := func() {
		called.Store(true)
	}

	w, err := Watch(path, callback, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Stop immediately
	w.Stop()

	// Try to modify file
	time.Sleep(100 * time.Millisecond)
	err = os.WriteFile(path, []byte("port: 9000"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)

	if called.Load() {
		t.Error("callback should not be called after Stop()")
	}
}

func TestWatch_FileNotFound(t *testing.T) {
	path := "/nonexistent/config.yaml"

	callback := func() {}

	_, err := Watch(path, callback, nil)
	if err == nil {
		t.Fatal("expected error when watching nonexistent file")
	}
}

func TestWatch_MultipleChanges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	err := os.WriteFile(path, []byte("port: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var count atomic.Int32
	callback := func() {
		count.Add(1)
	}

	w, err := Watch(path, callback, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	time.Sleep(100 * time.Millisecond)

	// Make multiple changes
	for i := 0; i < 3; i++ {
		err = os.WriteFile(path, []byte("port: 9000"), 0644)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(200 * time.Millisecond)
	}

	if count.Load() == 0 {
		t.Error("callback should be called at least once")
	}
}
