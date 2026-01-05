package circuit_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"github.com/moq77111113/circuit"
)

// ExampleFrom_minimal demonstrates the simplest Circuit integration.
// Creates a handler for a config struct and mounts it on a standard http.ServeMux.
func ExampleFrom_minimal() {
	// Define config struct with circuit tags
	type Config struct {
		Host string `yaml:"host" circuit:"type:text,help:Server hostname"`
		Port int    `yaml:"port" circuit:"type:number,help:Server port,required,min:1,max:65535"`
		TLS  bool   `yaml:"tls" circuit:"type:checkbox,help:Enable TLS"`
	}

	// Create temporary config file
	tmpDir := os.TempDir()
	configPath := filepath.Join(tmpDir, "example_minimal.yaml")
	_ = os.WriteFile(configPath, []byte("host: localhost\nport: 8080\ntls: false\n"), 0644)
	defer os.Remove(configPath)

	var cfg Config

	// Create Circuit handler
	handler, err := circuit.From(&cfg,
		circuit.WithPath(configPath),
		circuit.WithTitle("Example Config"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Mount on standard mux
	mux := http.NewServeMux()
	mux.Handle("/config", handler)

	// Test with httptest
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, _ := http.Get(server.URL + "/config")
	fmt.Println("Status:", resp.StatusCode)

	// Output:
	// Status: 200
}

// ExampleFrom_auth demonstrates Basic Auth integration.
// Shows both plaintext (dev) and argon2id (production) authentication.
func ExampleFrom_auth() {
	type Config struct {
		Host string `yaml:"host" circuit:"type:text"`
		Port int    `yaml:"port" circuit:"type:number"`
	}

	// Create temporary config file
	tmpDir := os.TempDir()
	configPath := filepath.Join(tmpDir, "example_auth.yaml")
	_ = os.WriteFile(configPath, []byte("host: localhost\nport: 8080\n"), 0644)
	defer os.Remove(configPath)

	var cfg Config

	// Development: plaintext password (DO NOT use in production)
	// auth := circuit.NewBasicAuth("admin", "dev-password")

	// Production: argon2id PHC hash
	// Hash generated with: golang.org/x/crypto/argon2
	// This example uses a plaintext password for testing purposes
	auth := circuit.NewBasicAuth("admin", "secure-password")

	handler, err := circuit.From(&cfg,
		circuit.WithPath(configPath),
		circuit.WithAuth(auth),
	)
	if err != nil {
		log.Fatal(err)
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	// Request without auth - should fail
	resp, _ := http.Get(server.URL)
	fmt.Println("Without auth:", resp.StatusCode)

	// Request with auth - should succeed
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.SetBasicAuth("admin", "secure-password")
	resp, _ = http.DefaultClient.Do(req)
	fmt.Println("With auth:", resp.StatusCode)

	// Output:
	// Without auth: 401
	// With auth: 200
}

// ExampleFrom_actions demonstrates executable actions.
// Actions appear as buttons in the UI and run server-side code.
func ExampleFrom_actions() {
	type Config struct {
		CacheSize int `yaml:"cache_size" circuit:"type:number"`
	}

	// Create temporary config file
	tmpDir := os.TempDir()
	configPath := filepath.Join(tmpDir, "example_actions.yaml")
	_ = os.WriteFile(configPath, []byte("cache_size: 100\n"), 0644)
	defer os.Remove(configPath)

	var cfg Config

	// Define safe actions with timeouts
	// NOTE: This is a dummy action for demonstration
	flushCache := circuit.NewAction("flush_cache", "Flush Cache", func(ctx context.Context) error {
		// In real code, this would clear your cache
		time.Sleep(100 * time.Millisecond) // Simulate work
		return nil
	}).Describe("Clears all cached data")

	// Destructive actions should require confirmation
	restart := circuit.NewAction("restart_worker", "Restart Worker", func(ctx context.Context) error {
		// In real code, this would restart your worker
		time.Sleep(100 * time.Millisecond) // Simulate work
		return nil
	}).Describe("Safely restarts the background worker").Confirm().WithTimeout(10 * time.Second)

	handler, err := circuit.From(&cfg,
		circuit.WithPath(configPath),
		circuit.WithActions(flushCache, restart),
	)
	if err != nil {
		log.Fatal(err)
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, _ := http.Get(server.URL)
	fmt.Println("Status:", resp.StatusCode)

	// Output:
	// Status: 200
}

// ExampleFrom_onChange demonstrates config change notifications.
// Circuit notifies your app when config changes, allowing you to apply updates.
func ExampleFrom_onChange() {
	type Config struct {
		Workers int `yaml:"workers" circuit:"type:number"`
	}

	tmpDir := os.TempDir()
	configPath := filepath.Join(tmpDir, "example_onchange.yaml")
	_ = os.WriteFile(configPath, []byte("workers: 4\n"), 0644)
	defer os.Remove(configPath)

	var cfg Config

	handler, err := circuit.From(&cfg,
		circuit.WithPath(configPath),
		circuit.WithOnChange(func(e circuit.ChangeEvent) {
			fmt.Printf("Config changed from %s\n", e.Source)
			// Your responsibility: apply new config to running components
			// Example: workerPool.Resize(cfg.Workers)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	fmt.Println("Handler created")

	// Output:
	// Handler created
}

// ExampleNewForwardAuth demonstrates reverse proxy authentication.
// Common with OAuth2 Proxy, Traefik ForwardAuth, Cloudflare Access.
func ExampleNewForwardAuth() {
	type Config struct {
		Setting string `yaml:"setting" circuit:"type:text"`
	}

	tmpDir := os.TempDir()
	configPath := filepath.Join(tmpDir, "example_forward.yaml")
	_ = os.WriteFile(configPath, []byte("setting: value\n"), 0644)
	defer os.Remove(configPath)

	var cfg Config

	// Configure forward auth to read headers from reverse proxy
	auth := circuit.NewForwardAuth("X-Forwarded-User", map[string]string{
		"email": "X-Forwarded-Email",
		"role":  "X-Auth-Role",
	})

	handler, err := circuit.From(&cfg,
		circuit.WithPath(configPath),
		circuit.WithAuth(auth),
	)
	if err != nil {
		log.Fatal(err)
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	// Request with proxy headers - should succeed
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("X-Forwarded-User", "user@example.com")
	req.Header.Set("X-Forwarded-Email", "user@example.com")
	resp, _ := http.DefaultClient.Do(req)
	fmt.Println("Status:", resp.StatusCode)

	// Output:
	// Status: 200
}
