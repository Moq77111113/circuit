package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/moq77111113/circuit"
)

type Config struct {
	Server ServerConfig `yaml:"server" circuit:"label:Server Configuration"`
	Cache  CacheConfig  `yaml:"cache" circuit:"label:Cache Configuration"`
}

type ServerConfig struct {
	Host string `yaml:"host" circuit:"type:text,label:Host,help:Server hostname"`
	Port int    `yaml:"port" circuit:"type:number,label:Port,help:Server port"`
}

type CacheConfig struct {
	Enabled bool   `yaml:"enabled" circuit:"type:checkbox,label:Enabled"`
	TTL     int    `yaml:"ttl" circuit:"type:number,label:TTL (seconds)"`
	Backend string `yaml:"backend" circuit:"type:select,options:redis|memcached|memory,label:Backend"`
}

func main() {
	cfg := Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Cache: CacheConfig{
			Enabled: true,
			TTL:     300,
			Backend: "memory",
		},
	}

	ui, err := circuit.From(&cfg,
		circuit.WithPath("config.yaml"),
		circuit.WithTitle("Circuit Actions Demo"),
		circuit.WithActions(
			circuit.Action{
				Name:  "restart",
				Label: "Restart Service",
				Run: func(ctx context.Context) error {
					log.Println("ACTION: Restarting service...")
					time.Sleep(2 * time.Second)
					log.Println("ACTION: Service restarted successfully")
					return nil
				},
			}.Describe("Restart the application service (requires confirmation)").Confirm().Timeout(30 * time.Second),

			circuit.Action{
				Name:  "reload",
				Label: "Reload Configuration",
				Run: func(ctx context.Context) error {
					log.Println("ACTION: Reloading configuration...")
					time.Sleep(1 * time.Second)
					log.Printf("ACTION: Configuration reloaded (Host=%s, Port=%d)", cfg.Server.Host, cfg.Server.Port)
					return nil
				},
			}.Describe("Hot reload configuration without downtime").Timeout(10 * time.Second),

			circuit.Action{
				Name:  "health-check",
				Label: "Health Check",
				Run: func(ctx context.Context) error {
					log.Println("ACTION: Running health check...")
					time.Sleep(500 * time.Millisecond)

					// Simulate health check
					if cfg.Server.Port == 0 {
						return fmt.Errorf("invalid port configuration")
					}

					log.Printf("ACTION: Health check passed (Server listening on %s:%d)", cfg.Server.Host, cfg.Server.Port)
					return nil
				},
			}.Describe("Verify service health and connectivity").Timeout(5 * time.Second),

			circuit.Action{
				Name:  "clear-cache",
				Label: "Clear Cache",
				Run: func(ctx context.Context) error {
					if !cfg.Cache.Enabled {
						return fmt.Errorf("cache is not enabled")
					}

					log.Printf("ACTION: Clearing %s cache...", cfg.Cache.Backend)
					time.Sleep(1 * time.Second)
					log.Println("ACTION: Cache cleared successfully")
					return nil
				},
			}.Describe("Clear all cached data").Timeout(15 * time.Second),

			circuit.Action{
				Name:  "fail-example",
				Label: "Failing Action (Demo)",
				Run: func(ctx context.Context) error {
					log.Println("ACTION: Running failing action...")
					time.Sleep(500 * time.Millisecond)
					return fmt.Errorf("this action intentionally fails for demo purposes")
				},
			}.Describe("Demonstrates error handling").Confirm(),

			circuit.Action{
				Name:  "git-status",
				Label: "Git Status",
				Run: func(ctx context.Context) error {
					cmd := exec.CommandContext(ctx, "git", "status", "--short")
					output, err := cmd.CombinedOutput()
					if err != nil {
						return fmt.Errorf("git command failed: %w", err)
					}
					log.Printf("ACTION: Git status:\n%s", string(output))
					return nil
				},
			}.Describe("Show current git repository status"),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting Circuit Actions Demo on http://%s", addr)
	log.Printf("Try the following actions:")
	log.Printf("  - Health Check: verify service configuration")
	log.Printf("  - Reload: hot reload configuration")
	log.Printf("  - Clear Cache: clear cache (only works if cache enabled)")
	log.Printf("  - Restart: restart service (requires confirmation)")
	log.Printf("  - Failing Action: demonstrates error handling")
	log.Printf("  - Git Status: shows git repository status")

	if err := http.ListenAndServe(addr, ui); err != nil {
		log.Fatal(err)
	}
}
