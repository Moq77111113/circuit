package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/moq77111113/circuit"
)

type Config struct {
	Host     string `yaml:"host" circuit:"type:text,help:Server hostname"`
	Port     int    `yaml:"port" circuit:"type:number,help:Server port,required"`
	LogLevel string `yaml:"log_level" circuit:"type:select,options:debug|info|warn|error,help:Logging level"`
}

func main() {
	// Parse flags
	authMode := flag.String("auth", "none", "Auth mode: none, basic, or forward")
	username := flag.String("user", "admin", "Username for basic auth")
	password := flag.String("pass", "secret", "Password for basic auth (plaintext or argon2id hash)")
	flag.Parse()

	var cfg Config

	// Configure authentication
	var auth circuit.Authenticator
	switch *authMode {
	case "basic":
		// Basic Auth: simple username + password
		// Password can be plaintext (dev) or argon2id hash (prod)
		auth = circuit.NewBasicAuth(*username, *password)
		log.Printf("Using Basic Auth (user: %s)", *username)

	case "forward":
		// Forward Auth: expects X-Forwarded-User header from reverse proxy
		// Use this when behind OAuth2 Proxy, Traefik, or Cloudflare Access
		auth = circuit.NewForwardAuth("X-Forwarded-User", map[string]string{
			"email": "X-Forwarded-Email",
		})
		log.Println("Using Forward Auth (expecting X-Forwarded-User header)")

	case "none":
		// No auth: open access
		auth = nil
		log.Println("⚠️  No authentication - UI is open to anyone!")

	default:
		log.Fatalf("Invalid auth mode: %s (use: none, basic, or forward)", *authMode)
	}

	handler, err := circuit.From(&cfg,
		circuit.WithPath("config.yaml"),
		circuit.WithTitle("Authenticated Settings"),
		circuit.WithAuth(auth), // Remove this line for no auth
		circuit.WithOnChange(func(e circuit.ChangeEvent) {
			log.Printf("Config changed: %v", e.Source)
			log.Printf("New values: Host=%s Port=%d LogLevel=%s", cfg.Host, cfg.Port, cfg.LogLevel)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Auth example running on :%s", 8080)
	log.Printf("Try accessing: http://localhost:%s", 8080)
	switch *authMode {
	case "basic":
		log.Printf("  Login with: %s / %s", *username, maskPassword(*password))
	case "forward":
		log.Println("  Must have X-Forwarded-User header set")
	}

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func maskPassword(pass string) string {
	if len(pass) <= 8 {
		return "***"
	}
	return pass[:4] + "***"
}
