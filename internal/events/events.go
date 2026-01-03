package events

// Source indicates where a configuration change originated.
type Source string

const (
	SourceFormSubmit Source = "form_submit"
	SourceFileChange Source = "file_change"
	SourceManual     Source = "manual"
)

// ChangeEvent describes a configuration change.
type ChangeEvent struct {
	Source Source
	Path   string
}

// OnChange is called when configuration changes.
type OnChange func(ChangeEvent)
