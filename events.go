package circuit

import "github.com/moq77111113/circuit/internal/sync"

type Source = sync.Source

const (
	SourceFormSubmit = sync.SourceFormSubmit
	SourceFileChange = sync.SourceFileChange
	SourceManual     = sync.SourceManual
)

type ChangeEvent = sync.ChangeEvent

type OnChange = sync.OnChange
