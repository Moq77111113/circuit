package reload

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestLoad_InitialLoad(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	err := os.WriteFile(path, []byte("port: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var cfg Config
	loader, err := Load(path, &cfg, func() {})
	if err != nil {
		t.Fatal(err)
	}
	defer loader.Stop()

	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
}

func TestLoad_ReloadOnChange(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	err := os.WriteFile(path, []byte("port: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var cfg Config
	var callbackCalled atomic.Bool

	loader, err := Load(path, &cfg, func() {
		callbackCalled.Store(true)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer loader.Stop()

	var port int
	loader.WithLock(func() {
		port = cfg.Port
	})

	if port != 8080 {
		t.Errorf("initial port should be 8080, got %d", port)
	}

	// Wait for watcher to start
	time.Sleep(100 * time.Millisecond)

	// Modify file
	err = os.WriteFile(path, []byte("port: 9000"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for reload
	time.Sleep(300 * time.Millisecond)

	loader.WithLock(func() {
		port = cfg.Port
	})

	if port != 9000 {
		t.Errorf("expected port 9000 after reload, got %d", port)
	}

	if !callbackCalled.Load() {
		t.Error("callback should be called after reload")
	}
}

func TestLoad_InvalidFile(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	path := "/nonexistent/config.yaml"
	var cfg Config

	_, err := Load(path, &cfg, func() {})
	if err == nil {
		t.Fatal("expected error for invalid file path")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	err := os.WriteFile(path, []byte("port: [invalid"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var cfg Config
	_, err = Load(path, &cfg, func() {})
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoad_Stop(t *testing.T) {
	type Config struct {
		Port int `yaml:"port"`
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	err := os.WriteFile(path, []byte("port: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var cfg Config
	var callbackCalled atomic.Bool

	loader, err := Load(path, &cfg, func() {
		callbackCalled.Store(true)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Stop immediately
	loader.Stop()

	time.Sleep(100 * time.Millisecond)

	// Modify file
	err = os.WriteFile(path, []byte("port: 9000"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)

	var port int
	loader.WithLock(func() {
		port = cfg.Port
	})

	if port != 8080 {
		t.Errorf("port should remain 8080 after stop, got %d", port)
	}

	if callbackCalled.Load() {
		t.Error("callback should not be called after Stop()")
	}
}

func TestLoad_MultipleReloads(t *testing.T) {
	type Config struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	err := os.WriteFile(path, []byte("host: localhost\nport: 8080"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var cfg Config
	var count atomic.Int32

	loader, err := Load(path, &cfg, func() {
		count.Add(1)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer loader.Stop()

	time.Sleep(100 * time.Millisecond)

	// First change
	err = os.WriteFile(path, []byte("host: example.com\nport: 9000"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)

	var host string
	var port int
	loader.WithLock(func() {
		host = cfg.Host
		port = cfg.Port
	})

	if host != "example.com" || port != 9000 {
		t.Errorf("expected example.com:9000, got %s:%d", host, port)
	}

	// Second change
	err = os.WriteFile(path, []byte("host: test.com\nport: 3000"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)

	loader.WithLock(func() {
		host = cfg.Host
		port = cfg.Port
	})

	if host != "test.com" || port != 3000 {
		t.Errorf("expected test.com:3000, got %s:%d", host, port)
	}

	if count.Load() < 2 {
		t.Errorf("expected at least 2 callbacks, got %d", count.Load())
	}
}
