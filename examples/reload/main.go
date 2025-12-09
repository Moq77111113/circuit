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

	// Setup Circuit UI
	handler, err := circuit.From(&app.config,
		circuit.WithPath("config.yaml"),
		circuit.WithTitle("Reload Example"),
		circuit.WithOnChange(func(e circuit.ChangeEvent) {
			app.UpdateConfig(app.config)
		}),
	)
	if err != nil {
		panic(err)
	}

	println("Reload example running on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
