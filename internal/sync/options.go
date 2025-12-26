package sync

type Option func(*Store)

// WithAutoApply controls whether POST automatically updates memory
func WithAutoApply(enable bool) Option {
	return func(s *Store) {
		s.autoApply = enable
	}
}

// WithAutoSave controls whether changes automatically persist to disk.
func WithAutoSave(enable bool) Option {
	return func(s *Store) {
		s.autoSave = enable
	}
}

// WithSaveFunc sets a custom function to persist config changes.
func WithSaveFunc(fn SaveFunc) Option {
	return func(s *Store) {
		s.saveFunc = fn
	}
}

// WithOnChange registers a callback for configuration change events.
func WithOnChange(fn OnChange) Option {
	return func(s *Store) {
		s.onChange = fn
	}
}

func WithOnError(fn func(error)) Option {
	return func(s *Store) {
		s.onError = fn
	}
}
