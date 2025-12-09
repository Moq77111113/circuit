package circuit

import "github.com/moq77111113/circuit/internal/reload"

type Source = reload.Source

const (
	SourceFormSubmit = reload.SourceFormSubmit
	SourceFileChange = reload.SourceFileChange
	SourceManual     = reload.SourceManual
)

type ChangeEvent = reload.ChangeEvent

type OnChange = reload.OnChange
