package circuit

// SaveFunc is called to persist configuration changes.
// Receives the current config value and path, returns error if persistence fails.
type SaveFunc func(cfg any, path string) error

// Option configures behavior passed to `From`.
type Option func(*config)

type config struct {
	path          string
	title         string
	brand         bool
	readOnly      bool
	onChange      OnChange
	onError       func(error)
	autoReload    bool
	autoApply     bool
	autoSave      bool
	saveFunc      SaveFunc
	authenticator Authenticator
	actions       []Action
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

// WithOnError registers a callback for errors during file watch or auto-reload.
// Useful for logging or alerting when config updates fail.
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
// Use handler.Apply(formData) to manually confirm changes after review.
func WithAutoApply(enable bool) Option {
	return func(c *config) {
		c.autoApply = enable
	}
}

// WithAutoSave controls whether changes trigger disk persistence.
// If false: memory update happens, but no disk save.
// Use handler.Save() to manually persist changes when ready.
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

// WithReadOnly controls whether the UI is read-only.
// When true, all inputs are disabled and Save/Add/Remove buttons are hidden.
func WithReadOnly(enable bool) Option {
	return func(c *config) {
		c.readOnly = enable
	}
}

// WithActions registers executable actions in the UI.
// Actions are displayed in the Actions section and can be triggered by operators.
func WithActions(actions ...Action) Option {
	return func(c *config) {
		c.actions = actions
	}
}
