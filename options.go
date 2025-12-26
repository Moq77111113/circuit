package circuit

import (
	"github.com/moq77111113/circuit/internal/auth"
	"github.com/moq77111113/circuit/internal/sync"
)

// SaveFunc is called to persist config changes.s
type SaveFunc = sync.SaveFunc

// Option configures behavior passed to `From`.
type Option func(*config)

type config struct {
	path          string
	title         string
	brand         bool
	onChange      OnChange
	onError       func(error)
	autoReload    bool
	autoApply     bool
	autoSave      bool
	saveFunc      SaveFunc
	authenticator auth.Authenticator
}

// WithPath sets the filesystem path to the YAML configuration file that the
// loader will read and watch. This option is required for `From` to succeed.
func WithPath(path string) Option {
	return func(c *config) {
		c.path = path
	}
}

// WithTitle sets the title displayed in the UI page.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}

// WithBrand controls whether the Circuit footer/brand is shown in the UI.
// The default is true.
func WithBrand(b bool) Option {
	return func(c *config) {
		c.brand = b
	}
}

// WithOnChange registers a callback for configuration change events.
// The callback receives a ChangeEvent indicating the source of the change.
func WithOnChange(fn OnChange) Option {
	return func(c *config) {
		c.onChange = fn
	}
}

func WithOnError(fn func(error)) Option {
	return func(c *config) {
		c.onError = fn
	}
}

// WithAutoWatch controls whether file watching is enabled.
// When true (default), changes to the YAML file trigger automatic reload.
// When false, file watching is disabled and reloads must be manual.
func WithAutoWatch(enable bool) Option {
	return func(c *config) {
		c.autoReload = enable
	}
}

// WithAuth sets the authenticator for the Circuit UI.
// If not provided, the UI is accessible without authentication.
func WithAuth(a Authenticator) Option {
	return func(c *config) {
		c.authenticator = a
	}
}

// WithAutoApply controls whether POST automatically updates in-memory config.
// If false: POST renders preview form with submitted values (doesn't modify memory).
func WithAutoApply(enable bool) Option {
	return func(c *config) {
		c.autoApply = enable
	}
}

// WithAutoSave controls whether changes trigger disk persistence.
// If false: memory update happens, but no disk save. User must call Save() manually.
func WithAutoSave(enable bool) Option {
	return func(c *config) {
		c.autoSave = enable
	}
}

// WithSaveFunc allows custom persistence implementation.
// Replaces default file writing with custom logic.
func WithSaveFunc(fn SaveFunc) Option {
	return func(c *config) {
		c.saveFunc = fn
	}
}
