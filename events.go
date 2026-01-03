package circuit

import "github.com/moq77111113/circuit/internal/events"

// Source indicates where a configuration change originated.
type Source = events.Source

const (
	SourceFormSubmit = events.SourceFormSubmit
	SourceFileChange = events.SourceFileChange
	SourceManual     = events.SourceManual
)

// ChangeEvent describes a configuration change.
type ChangeEvent = events.ChangeEvent

// OnChange is called when configuration changes.
type OnChange = events.OnChange
