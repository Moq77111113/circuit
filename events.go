package circuit

import "github.com/moq77111113/circuit/internal/events"

// Source indicates where a configuration change originated.
//
// Use the Source field in ChangeEvent to determine the origin of a config update.
type Source = events.Source

const (
	// SourceFormSubmit indicates the change came from a user submitting the web form.
	SourceFormSubmit = events.SourceFormSubmit

	// SourceFileChange indicates the change came from the file watcher detecting
	// an external modification to the config file.
	SourceFileChange = events.SourceFileChange

	// SourceManual indicates the change came from calling handler.Apply() directly.
	SourceManual = events.SourceManual
)

// ChangeEvent describes a configuration change.
//
// Delivered to OnChange callbacks after the in-memory config has been updated.
// Your application is responsible for applying the new config to running components.
type ChangeEvent = events.ChangeEvent

// OnChange is called when configuration changes.
//
// The callback is invoked AFTER the in-memory config struct has been updated.
// Use this to:
//   - Apply new config to running services
//   - Log config changes
//   - Trigger restart of affected components
//   - Validate config before applying
//
// Example:
//
//	circuit.WithOnChange(func(e circuit.ChangeEvent) {
//	    log.Printf("config updated from %s", e.Source)
//	    server.ApplyConfig(cfg)
//	})
type OnChange = events.OnChange
