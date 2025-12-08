package circuit

// Option configures behavior passed to `From`.
type Option func(*config)

type config struct {
	path    string
	title   string
	brand   bool
	onApply func()
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

// OnApply registers a callback that is invoked after a successful reload of
// the configuration (for example, when the watched YAML file changes and the
// new values are parsed). The callback runs asynchronously from the watcher.
func OnApply(fn func()) Option {
	return func(c *config) {
		c.onApply = fn
	}
}

// WithBrand controls whether the Circuit footer/brand is shown in the UI.
// The default is true.
func WithBrand(b bool) Option {
	return func(c *config) {
		c.brand = b
	}
}
