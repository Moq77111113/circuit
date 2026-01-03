package sync

import "github.com/moq77111113/circuit/internal/events"

// Type aliases for internal convenience.
type Source = events.Source
type ChangeEvent = events.ChangeEvent
type OnChange = events.OnChange

const (
	SourceFormSubmit = events.SourceFormSubmit
	SourceFileChange = events.SourceFileChange
	SourceManual     = events.SourceManual
)
