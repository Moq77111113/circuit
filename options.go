package circuit

// Option configures the circuit UI.
type Option func(*config)

type config struct {
	path    string
	title   string
	brand   bool
	onApply func()
}

// WithPath sets the path to the YAML config file.
func WithPath(path string) Option {
	return func(c *config) {
		c.path = path
	}
}

// WithTitle sets a custom title for the UI page.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}

// OnApply sets a callback to be called after config changes.
func OnApply(fn func()) Option {
	return func(c *config) {
		c.onApply = fn
	}
}

// WithBrand sets whether to show the Circuit brand in the UI footer. (default: true)
func WithBrand(b bool) Option {
	return func(c *config) {
		c.brand = b
	}
}
