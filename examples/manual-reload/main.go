package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/moq77111113/circuit"
)

type AppConfig struct {
	Name    string   `yaml:"name" circuit:"type:text,help:Application name"`
	Version string   `yaml:"version" circuit:"type:text,help:Version number"`
	Tags    []string `yaml:"tags" circuit:"type:text,help:Application tags"`
}

func main() {
	var cfg AppConfig

	handler, err := circuit.From(&cfg,
		circuit.WithPath("config.yaml"),
		circuit.WithTitle("Manual Reload Example"),
		circuit.WithAutoReload(false),
		circuit.WithOnChange(func(e circuit.ChangeEvent) {
			source := map[circuit.Source]string{
				circuit.SourceFormSubmit: "form submit",
				circuit.SourceFileChange: "file change",
				circuit.SourceManual:     "manual reload",
			}[e.Source]
			log.Printf("[%s] Config updated: %s v%s (tags: %v)\n",
				source, cfg.Name, cfg.Version, cfg.Tags)
		}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Manual Reload Example running on http://localhost:8081")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  - Auto-reload disabled (WithAutoReload(false))")
	fmt.Println("  - File changes won't trigger reload automatically")
	fmt.Println("  - Only form submits and manual reloads work")
	fmt.Println()
	fmt.Println("Try:")
	fmt.Println("  1. Edit config.yaml manually")
	fmt.Println("  2. Notice it doesn't reload automatically")
	fmt.Println("  3. Submit form to trigger reload")
	fmt.Println()

	if err := http.ListenAndServe(":8081", handler); err != nil {
		panic(err)
	}
}
