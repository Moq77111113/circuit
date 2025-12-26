package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/moq77111113/circuit"
)

type Config struct {
	Message     string `yaml:"message" circuit:"type:text,help:Message to display"`
	RefreshRate int    `yaml:"refresh_rate" circuit:"type:number,help:Refresh rate in seconds,min:1"`
}

type App struct {
	mu     sync.Mutex
	config Config
}

func (a *App) Run() {
	for {
		a.mu.Lock()
		rate := a.config.RefreshRate
		msg := a.config.Message
		a.mu.Unlock()

		if rate == 0 {
			rate = 5
		}

		fmt.Printf("[%s] App running: %s\n", time.Now().Format(time.TimeOnly), msg)
		time.Sleep(time.Duration(rate) * time.Second)
	}
}

func (a *App) UpdateConfig(newCfg Config) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.config = newCfg
	fmt.Println(">>> Configuration updated!")
}

func main() {
	app := &App{
		config: Config{
			Message:     "Initial message",
			RefreshRate: 5,
		},
	}

	// Start the "app" in background
	go app.Run()

	// Setup Circuit UI with preview mode
	handler, err := circuit.From(&app.config,
		circuit.WithPath("config.yaml"),
		circuit.WithTitle("Preview Mode Example"),
		circuit.WithAutoApply(false),
		circuit.WithOnChange(func(e circuit.ChangeEvent) {
			app.UpdateConfig(app.config)
		}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Preview Mode example running on :8080")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  - Preview mode enabled (WithAutoApply(false))")
	fmt.Println("  - Changes require confirmation before being applied")
	fmt.Println("  - Watch the app output to see when config actually updates")
	fmt.Println()
	fmt.Println("Try:")
	fmt.Println("  1. Edit a field and click Save")
	fmt.Println("  2. You'll see a preview with Confirm/Cancel buttons")
	fmt.Println("  3. Notice the app still uses the old config (check console output)")
	fmt.Println("  4. Click Confirm to apply changes")
	fmt.Println("  5. Now the app receives the updated config")
	fmt.Println()

	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
