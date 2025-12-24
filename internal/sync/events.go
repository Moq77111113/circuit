package sync

// Source represents the origin of a config change.
type Source uint8

const (
	SourceFormSubmit Source = iota
	SourceFileChange
	SourceManual
)

// ChangeEvent describes a config change that occurred.
type ChangeEvent struct {
	Source Source
	Path   string
}

// OnChange is called when the config changes.
type OnChange func(ChangeEvent)
