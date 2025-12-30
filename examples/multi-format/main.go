package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/moq77111113/circuit"
)

type ServerConfig struct {
	Port    int    `circuit:"input,number,help:Server port" yaml:"port" toml:"port" json:"port"`
	Host    string `circuit:"input,text,help:Server host" yaml:"host" toml:"host" json:"host"`
	Enabled bool   `circuit:"input,checkbox,help:Enable server" yaml:"enabled" toml:"enabled" json:"enabled"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: multi-format <config-file>")
		fmt.Println("Examples:")
		fmt.Println("  multi-format config.yaml")
		fmt.Println("  multi-format config.toml")
		fmt.Println("  multi-format config.json")
		os.Exit(1)
	}

	configPath := os.Args[1]

	cfg := ServerConfig{}
	ui, err := circuit.From(&cfg,
		circuit.WithPath(configPath),
		circuit.WithTitle("Multi-Format Config"),
	)
	if err != nil {
		log.Fatalf("Failed to create UI: %v", err)
	}

	fmt.Printf("Configuration format detected from: %s\n", configPath)
	fmt.Printf("Starting server on http://localhost:8080\n")
	fmt.Printf("Current config: port=%d host=%s enabled=%v\n", cfg.Port, cfg.Host, cfg.Enabled)

	log.Fatal(http.ListenAndServe(":8080", ui))
}
