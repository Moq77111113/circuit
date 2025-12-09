package circuit

// Option configures behavior passed to `From`.
type Option func(*config)

type config struct {
	path       string
	title      string
	brand      bool
	onChange   OnChange
	autoReload bool
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

// WithAutoReload controls whether file watching is enabled.
// When true (default), changes to the YAML file trigger automatic reload.
// When false, file watching is disabled and reloads must be manual.
func WithAutoReload(enable bool) Option {
	return func(c *config) {
		c.autoReload = enable
	}
}
