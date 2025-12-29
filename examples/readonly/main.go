package main

import (
	"log"
	"net/http"

	"github.com/moq77111113/circuit"
)

type Config struct {
	APIKey  string `circuit:"text,required,help:Your API key" yaml:"api_key"`
	Version string `circuit:"text,readonly,help:Application version (read-only)" yaml:"version"`
	Port    int    `circuit:"number,min:1,max:65535,help:Server port" yaml:"port"`
}

func main() {
	cfg := &Config{
		APIKey:  "sk-1234567890",
		Version: "v1.0.0",
		Port:    8080,
	}

	ui, err := circuit.From(cfg,
		circuit.WithPath("config.yaml"),
		circuit.WithTitle("Readonly Example"),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server running on http://localhost:8080")
	log.Println("Try editing:")
	log.Println("  - APIKey (editable)")
	log.Println("  - Version (readonly via tag)")
	log.Println("  - Port (editable)")

	if err := http.ListenAndServe(":8080", ui); err != nil {
		log.Fatal(err)
	}
}
